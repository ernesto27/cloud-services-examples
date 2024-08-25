package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const intervalServerError = "Internal server error"

func handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {

	println("DB connection successful")
	printDebug(event)
	var response events.APIGatewayProxyResponse

	switch event.RouteKey {
	case "GET /users":
		users, err := db.GetUsers()
		if err != nil {
			fmt.Println(err)
			return responseEndpoint(http.StatusInternalServerError, intervalServerError, ""), nil

		}

		usersJSON, err := json.Marshal(users)
		if err != nil {
			fmt.Println(err)
			return responseEndpoint(http.StatusInternalServerError, intervalServerError, ""), nil
		}

		response = responseEndpoint(http.StatusOK, string(usersJSON), "application/json")

	case "POST /users":
		var user User
		err := json.Unmarshal([]byte(event.Body), &user)
		if err != nil {
			fmt.Println(err)
			return responseEndpoint(http.StatusInternalServerError, intervalServerError, ""), nil
		}

		if user.Username == "" || user.Email == "" || user.PasswordFromPayload == "" {
			return responseEndpoint(http.StatusBadRequest, "username, email, and password are required", ""), nil
		}

		fmt.Println("USER: ", user)

		err = db.CreateUser(user)
		if err != nil {
			fmt.Println(err)
			return responseEndpoint(http.StatusInternalServerError, intervalServerError, ""), nil
		}

		return responseEndpoint(http.StatusCreated, "User created", ""), nil

	case "GET /users/{id}":
		id, err := strconv.Atoi(event.PathParameters["id"])
		if err != nil {
			fmt.Println(err)
			return responseEndpoint(http.StatusBadRequest, "Invalid ID", ""), nil
		}

		user, err := db.GetUserByID(id)
		if err != nil {
			fmt.Println(err)
			return responseEndpoint(http.StatusInternalServerError, intervalServerError, ""), nil
		}

		userJSON, err := json.Marshal(user)
		if err != nil {
			fmt.Println(err)
			return responseEndpoint(http.StatusInternalServerError, intervalServerError, ""), nil
		}

		response = responseEndpoint(http.StatusOK, string(userJSON), "application/json")

	case "PUT /users/{id}":
		id, err := strconv.Atoi(event.PathParameters["id"])
		if err != nil {
			fmt.Println(err)
			return responseEndpoint(http.StatusBadRequest, "Invalid ID", ""), nil
		}

		var user User
		err = json.Unmarshal([]byte(event.Body), &user)
		if err != nil {
			fmt.Println(err)
			return responseEndpoint(http.StatusInternalServerError, intervalServerError, ""), nil
		}

		if user.Username == "" || user.Email == "" || user.PasswordFromPayload == "" {
			return responseEndpoint(http.StatusBadRequest, "username, email, and password are required", ""), nil
		}

		user.ID = id

		err = db.UpdateUser(user)
		if err != nil {
			fmt.Println(err)
			return responseEndpoint(http.StatusInternalServerError, intervalServerError, ""), nil
		}

		response = responseEndpoint(http.StatusOK, "User updated", "")
	default:
		response = responseEndpoint(http.StatusNotFound, "Not found", "")
	}

	return response, nil
}

func responseEndpoint(statusCode int, message string, contentType string) events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       message,
	}

	if contentType != "" {
		resp.Headers = map[string]string{
			"Content-Type": contentType,
		}
	}

	return resp
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

var db *Mysql

func main() {
	var err error

	if err != nil {
		fmt.Println(err)
		panic(err)

	}
	lambda.Start(handler)
}
