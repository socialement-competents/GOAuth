package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/socialement-competents/goauth/models"
)

// Event : incoming event
type Event struct {
	ID int `json:"id"`
}

// HandleRequest : handle the incoming requests
func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	_ = models.User{}
	panic("NOT IMPLEMENTED")
}

func main() {
	lambda.Start(HandleRequest)
}
