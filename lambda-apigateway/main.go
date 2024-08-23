package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	printDebug(event)

	db, err := NewMysql(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE"))
	if err != nil {
		fmt.Println(err)

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Interval server error",
		}, nil
	}

	var response events.APIGatewayProxyResponse

	switch event.RouteKey {
	case "GET /users":
		users, err := db.GetUsers()
		if err != nil {
			fmt.Println(err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       "Interval server error",
			}, nil
		}

		usersJSON, err := json.Marshal(users)
		if err != nil {
			fmt.Println(err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       "Internal server error",
			}, nil
		}

		response = responseJSON(string(usersJSON))

	case "POST /users":
		response = events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "create user in DB",
		}

	case "GET /users/{id}":
		response = events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "get user by ID from DB",
		}

	case "PUT /users/{id}":
		response = events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "update user by ID in DB",
		}

	}

	return response, nil
}

func responseJSON(body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       body,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

func printDebug(event events.APIGatewayV2HTTPRequest) {
	fmt.Println("RouteKey", event.RouteKey)
	fmt.Println("RawPath", event.RawPath)

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
	fmt.Println("body", event.Body)
}

func main() {
	lambda.Start(handler)
}
