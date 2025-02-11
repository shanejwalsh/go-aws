package api

import (
	"fmt"
	"lambda-func/database"
	"lambda-func/types"
)

type ApiHandler struct {
	dbStore *database.DynamoDBClient
}

func NewApiHandler(dbStore *database.DynamoDBClient) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterUser(event types.RegisterUser) error {

	if (event.Username) == "" || event.Password == "" {
		return fmt.Errorf("username and password are required")
	}

	userExists, err := api.dbStore.DoesUserExist(event.Username)

	if err != nil {
		return fmt.Errorf("error checking if user exists %w", err)
	}

	if userExists {
		return fmt.Errorf("user already exists")
	}

	err = api.dbStore.InsertUser(event)

	if err != nil {
		return fmt.Errorf("error inserting user %w", err)
	}

	return nil
}
