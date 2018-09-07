package main

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

// Payload sent by the GitHub API
type Payload struct {
	// Required. The client ID you received from GitHub for your GitHub App.
	ClientID string `json:"client_id"`
	// Required. The client secret you received from GitHub for your GitHub App.
	ClientSecret string `json:"client_secret"`
	// Required. The code you received as a response to Step 1.
	Code string `json:"code"`
	// The URL in your application where users are sent after authorization.
	RedirectURI string `json:"redirect_uri"`
	// The unguessable random stdring you provided in Step 1.
	State string `json:"state"`
}

// RegisterUser : registers an user with GitHub
func RegisterUser(ctx context.Context, payload Payload) error {
	clientID := os.Getenv("GH_ID")
	clientSecret := os.Getenv("GH_SECRET")

	if clientID == "" {
		return errors.New("$GH_ID should be set")
	}
	if clientSecret == "" {
		return errors.New("$GH_SECRET should be set")
	}

	return nil
}

func main() {
	lambda.Start(RegisterUser)
}
