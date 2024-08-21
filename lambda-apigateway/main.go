package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("event", event)

	for k, v := range event.Headers {
		fmt.Printf("header %s: %s\n", k, v)
	}

	for k, v := range event.QueryStringParameters {
		fmt.Printf("query %s: %s\n", k, v)
	}

	for k, v := range event.PathParameters {
		fmt.Printf("path %s: %s\n", k, v)
	}

	for k, v := range event.StageVariables {
		fmt.Printf("stage %s: %s\n", k, v)
	}

	fmt.Println("requestContext", event.RequestContext)

	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "update all " + event.RouteKey,
	}
	return response, nil
}

func main() {
	lambda.Start(handler)
}
