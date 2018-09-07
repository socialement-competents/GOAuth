package main

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

type RegisterEvent struct {
}

// RegisterUser : registers an user with GitHub
func RegisterUser(ctx context.Context, evt RegisterEvent) error {
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
