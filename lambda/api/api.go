package api

import (
	"encoding/json"
	"fmt"
	"lambda-func/database"
	"lambda-func/types"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
	dbStore database.UserStore
}

func NewApiHandler(dbStore database.UserStore) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterUser(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var registerUserEvent types.RegisterUser

	err := json.Unmarshal([]byte(event.Body), &registerUserEvent)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "bad request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	if (registerUserEvent.Username) == "" || registerUserEvent.Password == "" {
		return events.APIGatewayProxyResponse{
			Body:       "bad request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	userExists, err := api.dbStore.DoesUserExist(registerUserEvent.Username)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "bad request",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	if userExists {
		return events.APIGatewayProxyResponse{
			Body:       "bad request - user already exists",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	user, err := types.NewUser(registerUserEvent)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	err = api.dbStore.InsertUser(user)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error inserting user %w", err)
	}

	return events.APIGatewayProxyResponse{
		Body:       event.Body,
		StatusCode: http.StatusCreated,
	}, nil
}

func (api ApiHandler) LoginUser(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var loginRequest LoginRequest

	err := json.Unmarshal([]byte(req.Body), &loginRequest)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "bad request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	user, err := api.dbStore.GetUser(loginRequest.Username)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "user with name " + loginRequest.Username + " not found",
			StatusCode: http.StatusNotFound,
		}, err
	}

	if !types.ValidatePassword(user.PasswordHash, loginRequest.Password) {
		return events.APIGatewayProxyResponse{
			Body:       "incorect password",
			StatusCode: http.StatusBadRequest,
		}, err

	}

	return events.APIGatewayProxyResponse{
		Body:       "Logged in",
		StatusCode: http.StatusOK,
	}, nil

}
