package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// NameEvent : payload
type NameEvent struct {
	Name string `json:"name"`
}

// HandleRequest : handle the incoming requests
func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	msg := fmt.Sprintf("Hello %s!", req.QueryStringParameters["name"])
	resp := events.APIGatewayProxyResponse{
		Body:       msg,
		StatusCode: 200,
	}
	return resp, nil
}

func main() {
	lambda.Start(HandleRequest)
}
