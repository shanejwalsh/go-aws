package database

import (
	"fmt"
	"lambda-func/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	TableName = "users"
)

type UserStore interface {
	DoesUserExist(username string) (bool, error)
	InsertUser(event types.User) error
	GetUser(username string) (types.User, error)
}

type DynamoDBClient struct {
	databaseStore *dynamodb.DynamoDB
}

func NewDynamoDBClient() DynamoDBClient {

	dbSession := session.Must(session.NewSession())
	db := dynamodb.New(dbSession)

	return DynamoDBClient{
		databaseStore: db,
	}
}

func (db DynamoDBClient) DoesUserExist(username string) (bool, error) {

	result, err := db.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})

	if err != nil {
		return true, err
	}

	if len(result.Item) > 0 {
		return true, nil
	}

	return false, nil
}

func (db DynamoDBClient) InsertUser(user types.User) error {

	item := map[string]*dynamodb.AttributeValue{
		"username": {
			S: aws.String(user.Username),
		},
		"password": {
			S: aws.String(user.PasswordHash),
		},
	}

	_, err := db.databaseStore.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item:      item,
	})

	if err != nil {
		return err
	}

	return nil
}

func (db DynamoDBClient) GetUser(username string) (types.User, error) {

	var user types.User

	result, err := db.databaseStore.GetItem(&dynamodb.GetItemInput{

		TableName: aws.String(TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})

	if err != nil {
		return user, err
	}

	if result.Item == nil {
		return user, fmt.Errorf("user with username %s not found", username)
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &user)

	if err != nil {
		return user, err

	}

	return user, nil
}
