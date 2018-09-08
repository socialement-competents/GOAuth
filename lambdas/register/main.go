package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// GHResponse : payload sent by the GitHub API
type GHResponse struct {
	Code string `json:"code"`
}

// GHToken : token given by GitHub in exchange for a Code
type GHToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// RegisterUser : registers an user with GitHub
func RegisterUser(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	clientID := os.Getenv("GH_ID")
	clientSecret := os.Getenv("GH_SECRET")

	if clientID == "" {
		return events.APIGatewayProxyResponse{
			Body:       "$GH_ID should be set",
			StatusCode: 400,
		}, nil
	}
	if clientSecret == "" {
		return events.APIGatewayProxyResponse{
			Body:       "$GH_SECRET should be set",
			StatusCode: 400,
		}, nil
	}

	var payload GHResponse

	data, err := json.Marshal(request.QueryStringParameters)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 400,
		}, err
	}

	err = json.Unmarshal(data, &payload)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 400,
		}, err
	}

	client := &http.Client{}

	url := fmt.Sprintf("https://github.com/login/oauth/access_token?code=%s&client_id=%s&client_secret=%s", payload.Code, clientID, clientSecret)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	fmt.Println("resp token: ", resp)

	if resp.StatusCode >= 400 {
		errtxt := fmt.Sprintf("bad GitHub response: %s", resp.Status)
		fmt.Println(errtxt)
		return events.APIGatewayProxyResponse{
			Body:       errtxt,
			StatusCode: resp.StatusCode,
		}, errors.New(errtxt)
	} else if resp.StatusCode >= 300 {
		errtxt := fmt.Sprintf("unexpected 3xx code: %s", resp.Status)
		fmt.Println(errtxt)
		return events.APIGatewayProxyResponse{
			Body:       errtxt,
			StatusCode: resp.StatusCode,
		}, errors.New(errtxt)
	}

	var token GHToken

	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 400,
		}, err
	}

	fmt.Println(token)

	req, err = http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("token %s", token.AccessToken))

	resp, err = client.Do(req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	fmt.Println("resp user: ", resp)

	if resp.StatusCode >= 400 {
		errtxt := fmt.Sprintf("bad GitHub response: %s", resp.Status)
		fmt.Println(errtxt)
		return events.APIGatewayProxyResponse{
			Body:       errtxt,
			StatusCode: resp.StatusCode,
		}, errors.New(errtxt)
	} else if resp.StatusCode >= 300 {
		errtxt := fmt.Sprintf("unexpected 3xx code: %s", resp.Status)
		fmt.Println(errtxt)
		return events.APIGatewayProxyResponse{
			Body:       errtxt,
			StatusCode: resp.StatusCode,
		}, errors.New(errtxt)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	user := buf.String()

	log.Println(user)
	fmt.Println(user)

	// check in the db if the user already exists
	// insert or update data

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       user,
	}, err
}

func main() {
	lambda.Start(RegisterUser)
}
