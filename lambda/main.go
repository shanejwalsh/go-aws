package main

import (
	"lambda-func/app"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {

	lambda.Start(app.NewApp().ApiHandler.RegisterUser)
}
