package main

import (
	"fmt"
	"lambda-func/app"
	"lambda-func/middleware"
	"lambda-func/types"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
)

func ProtetedHandler(req types.Req) (types.Res, error) {

	return types.Res{
		Body:       "secret",
		StatusCode: http.StatusOK,
	}, nil

}

func main() {

	appInstance := app.NewApp()

	handler := func(req types.Req) (types.Res, error) {
		switch req.Path {
		case "/register":
			return appInstance.ApiHandler.RegisterUser(req)

		case "/login":
			return app.NewApp().ApiHandler.LoginUser(req)
		case "/protected":
			return middleware.ValidateJWTMiddleware(ProtetedHandler)(req)
		default:
			return types.Res{
				Body:       "Not Found",
				StatusCode: http.StatusNotFound,
			}, fmt.Errorf("unknown path %s", req.Path)
		}

	}

	lambda.Start(handler)
}
