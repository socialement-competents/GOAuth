package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/socialement-competents/goauth/database"
	"github.com/socialement-competents/goauth/models"
)

// GHPayload : payload sent by the GitHub API
type GHPayload struct {
	Code string `json:"code"`
}

// GHToken : token given by GitHub in exchange for a Code
type GHToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

const accessTokenURL = "https://github.com/login/oauth/access_token?code=%s&client_id=%s&client_secret=%s"
const userURL = "https://api.github.com/user"
const githubProvider = "github"

var client = &http.Client{}
var clientID string
var clientSecret string

func init() {
	clientID = os.Getenv("GH_ID")
	clientSecret = os.Getenv("GH_SECRET")
}

// HandleCallback : handles the GitHub callback when calling "Connect with GitHub"
func HandleCallback(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if clientID == "" {
		return respond(http.StatusBadRequest, "$GH_ID should be set")
	}
	if clientSecret == "" {
		return respond(http.StatusBadRequest, "$GH_ID should be set")
	}

	payload, err := getCode(&request)
	if err != nil {
		return respond(
			http.StatusBadRequest,
			fmt.Sprintf("error getting the code from the payload: %v", err),
		)
	}

	token, err := getAcccessToken(payload)
	if err != nil {
		return respond(
			http.StatusInternalServerError,
			fmt.Sprintf("error getting the access token from GH: %v", err),
		)
	}

	ghuser, err := getUser(token)
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

	user, err := dbClient.GetUserByLogin(ghuser.Login, githubProvider)
	exists := true
	if err != nil {
		// the user doesn't exist, we will have to create it
		exists = false
		user = &models.User{Provider: githubProvider}
	}

	// Update our database with the newly fetched info
	// (we don't query api.github.com every time, because of rate limits)
	user.GHUser = ghuser
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

func getCode(request *events.APIGatewayProxyRequest) (*GHPayload, error) {
	var payload GHPayload
	data, err := json.Marshal(request.QueryStringParameters)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &payload)
	return &payload, err
}

func getAcccessToken(payload *GHPayload) (*GHToken, error) {
	url := fmt.Sprintf(
		accessTokenURL,
		payload.Code,
		clientID,
		clientSecret,
	)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if err = checkStatusCode(resp); err != nil {
		return nil, err
	}

	var token GHToken
	err = json.NewDecoder(resp.Body).Decode(&token)
	return &token, err
}

func getUser(token *GHToken) (*models.GHUser, error) {
	req, err := http.NewRequest(http.MethodGet, userURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("token %s", token.AccessToken))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if err = checkStatusCode(resp); err != nil {
		return nil, err
	}

	var user models.GHUser
	err = json.NewDecoder(resp.Body).Decode(&user)
	return &user, err
}

func checkStatusCode(resp *http.Response) error {
	if resp.StatusCode >= 400 {
		return fmt.Errorf("bad GitHub response: %s", resp.Status)
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
	lambda.Start(HandleCallback)
}
