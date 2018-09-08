package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

var client = &http.Client{}
var clientID string
var clientSecret string

func init() {
	clientID = os.Getenv("GH_ID")
	clientSecret = os.Getenv("GH_SECRET")
}

// RegisterUser : registers an user with GitHub
func RegisterUser(ctx context.Context, request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
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

	user, err := checkIfExists(ghuser)
	if err != nil {
		return respond(
			http.StatusInternalServerError,
			fmt.Sprintf("error check if the user exists: %v", err),
		)
	}

	log.Println(user)
	fmt.Println(user)

	// check in the db if the user already exists
	// insert or update data

	return respond(http.StatusOK, user)
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

func checkIfExists(user *models.GHUser) (*models.User, error) {
	return nil, errors.New("NOT IMPLEMENTED")
}

func respond(code int, body interface{}) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       fmt.Sprint(body),
	}
}

func main() {
	lambda.Start(RegisterUser)
}
