package hellolambda

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

// HelloWorldEvent : incoming event
type HelloWorldEvent struct {
	Name string `json:"name"`
}

// HandleRequest : handle the incoming requests
func HandleRequest(ctx context.Context, evt HelloWorldEvent) (string, error) {
	return fmt.Sprintf("Hello %s!", evt.Name), nil
}

func main() {
	lambda.Start(HandleRequest)
}
