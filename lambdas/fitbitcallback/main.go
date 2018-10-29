package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/socialement-competents/goauth/database"
	"github.com/socialement-competents/goauth/models"
)

// FitBitPayload : payload sent by the FitBit API
type FitBitPayload struct {
	UserID    string `json:"user_id"`
	ExpiresIn int    `json:"expires_in"`
	*FitBitToken
}

// FitBitToken : token given by FitBit in exchange for a Code
type FitBitToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

const userURL = "https://api.fitbit.com/1/user/%s/profile.json"

const githubProvider = "fitbit"

var client = &http.Client{}
var clientSecret string

func init() {
	clientSecret = os.Getenv("FITBIT_SECRET")
}

// FitBitCallback : handles the FitBit callback when calling "Connect with FitBit"
func FitBitCallback(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if clientSecret == "" {
		return respond(http.StatusBadRequest, "$FITBIT_SECRET should be set")
	}

	payload, err := parsePath(request.QueryStringParameters)
	if err != nil {
		return respond(
			http.StatusInternalServerError,
			fmt.Sprintf("error getting the access token from FitBit: %v", err),
		)
	}

	fitbituser, err := getUser(payload.FitBitToken, payload.UserID)
	if err != nil {
		return respond(
			http.StatusInternalServerError,
			fmt.Sprintf("error getting the user from GH: %v", err),
		)
	}

	dbClient, err := database.NewClient()
	if err != nil {
		return respond(
			http.StatusInternalServerError,
			fmt.Sprintf("couldn't connect to the db: %v", err.Error()),
		)
	}

	user, err := dbClient.GetUserByIdentifier("fitbit_id", fitbituser.EncodedID)
	exists := true
	if err != nil {
		// the user doesn't exist, we will have to create it
		exists = false
		user = &models.User{Provider: githubProvider}
	}

	// Update our database with the newly fetched info
	// (we don't query the FitBit API every time, because of rate limits)
	user.FitBitUser = fitbituser
	user.LastLogin = time.Now()

	var (
		verb       string
		statusCode int
	)

	if !exists {
		// the user wasn't previously found, we need to create it
		id, err := dbClient.CreateUser(user)
		if err != nil {
			return respond(
				http.StatusInternalServerError,
				fmt.Sprintf("creating the user failed: %v", err.Error()),
			)
		}
		user.ID = id
		verb = "created"
		statusCode = http.StatusCreated
	} else {
		// otherwise update it
		if err = dbClient.UpdateUser(user); err != nil {
			return respond(
				http.StatusInternalServerError,
				fmt.Sprintf("updating the user failed: %v", err.Error()),
			)
		}
		verb = "updated"
		statusCode = http.StatusOK
	}

	jsonBytes, err := json.Marshal(user)
	if err != nil {

		text := fmt.Sprintf("user %s but could not format to JSON: %v", verb, err)
		return respond(http.StatusAccepted, text)
	}

	return respond(statusCode, string(jsonBytes))
}

func parsePath(p map[string]string) (*FitBitPayload, error) {
	ttl, err := strconv.Atoi(p["expires_in"])
	if err != nil {
		log.Println("INVALID TTL: ", p["expires_in"])
		ttl = 575867
	}

	log.Print("Scope: " + p["scope"])

	fbp := &FitBitPayload{
		ExpiresIn: ttl,
		UserID:    p["user_id"],
		FitBitToken: &FitBitToken{
			AccessToken: p["access_token"],
			TokenType:   p["token_type"],
			Scope:       p["scope"],
		},
	}

	if !strings.Contains(fbp.Scope, "heartrate") {
		return nil, errors.New("the heartrate scope is required - scope: " + fbp.Scope)
	}

	return fbp, nil
}

func getUser(token *FitBitToken, userID string) (*models.FitBitUser, error) {
	url := fmt.Sprintf(userURL, userID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", token.TokenType, token.AccessToken))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if err = checkStatusCode(resp); err != nil {
		return nil, err
	}

	var user models.FitBitUser
	err = json.NewDecoder(resp.Body).Decode(&user)
	raw, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		user.RawPayload = string(raw)
	}
	return &user, err
}

func checkStatusCode(resp *http.Response) error {
	if resp.StatusCode >= 400 {
		return fmt.Errorf("bad FitBit response: %s", resp.Status)
	} else if resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected return code: %s", resp.Status)
	}

	return nil
}

func respond(code int, payload interface{}) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       fmt.Sprint(payload),
	}, nil
}

func main() {
	lambda.Start(FitBitCallback)
}
