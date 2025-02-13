package main

import (
	"fmt"
	"lambda-func/app"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {

	appInstance := app.NewApp()

	handler := func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.Path {
		case "/register":
			return appInstance.ApiHandler.RegisterUser(request)

		case "/login":
			return app.NewApp().ApiHandler.LoginUser(request)
		default:
			return events.APIGatewayProxyResponse{
				Body:       "Not Found",
				StatusCode: http.StatusNotFound,
			}, fmt.Errorf("unknown path %s", request.Path)
		}

	}

	lambda.Start(handler)
}
