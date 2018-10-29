package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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

// access_token=eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiIyMkQ0RE4iLCJzdWIiOiI1TEdIWUsiLCJpc3MiOiJGaXRiaXQiLCJ0eXAiOiJhY2Nlc3NfdG9rZW4iLCJzY29wZXMiOiJyaHIgcnBybyIsImV4cCI6MTU0MTQwNzIyNywiaWF0IjoxNTQwODAzNTQxfQ.biiNAuvU2FQsf39UbHimeK4amKkI7ARvKhUQoKrC2iQ
// user_id=5LGHYK
// scope=heartrate+profile
// token_type=Bearer
// expires_in=31536000

// REFRESH TOKEN : curl -X POST -i -H "Authorization: Basic MjJENEROOmMxN2NjMDczMjIzYzEyYmU4ZjUxNTk2OWI4NDcxYzIx" -H "Content-Type: application/x-www-form-urlencoded" -d "grant_type=refresh_token" -d "refresh_token=eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiIyMkQ0RE4iLCJzdWIiOiI1TEdIWUsiLCJpc3MiOiJGaXRiaXQiLCJ0eXAiOiJhY2Nlc3NfdG9rZW4iLCJzY29wZXMiOiJyaHIgcnBybyIsImV4cCI6MTU0MTQwNzIyNywiaWF0IjoxNTQwODAzNTQxfQ.biiNAuvU2FQsf39UbHimeK4amKkI7ARvKhUQoKrC2iQ" https://api.fitbit.com/oauth2/token

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

	payload, err := parseRequest(&request)
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

func parseRequest(request *events.APIGatewayProxyRequest) (*FitBitPayload, error) {
	var payload FitBitPayload
	data, err := json.Marshal(request.QueryStringParameters)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &payload)
	return &payload, err
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
