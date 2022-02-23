package user

import (
	"github.com/harrisjib216/Golang-Serverless-AWS/pkg/validators"

	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func GetUser(email, tableName string, dynamoClient dynamodbiface.DynamoDBAPI) (
	*User, error,
) {
	query := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	res, err := dynamoClient.GetItem(query)

	if err != nil {
		return nil, errors.New("Could not find this item")
	}

	item := new(User)
	err = dynamodbattribute.UnmarshalMap(res.Item, item)

	if err != nil {
		return nil, errors.New("Could convert this item")
	}

	return item, nil
}

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynamoClient dynamodbiface.DynamoDBAPI) (
	*User, error,
) {
	var newUser User

	if err := json.Unmarshal([]byte(req.Body), &newUser); err != nil {
		return nil, errors.New("Incorrect user data")
	}

	if !validators.IsEmailValid(newUser.Email) {
		return nil, errors.New("Email is invalid")
	}

	existingUser, _ := GetUser(newUser.Email, tableName, dynamoClient)
	if existingUser != nil && len(existingUser.Email) != 0 {
		return nil, errors.New("This user already exists")
	}

	dynamoMap, err := dynamodbattribute.MarshalMap(newUser)

	if err != nil {
		return nil, errors.New("Could not convert this item")
	}

	creationRequest := &dynamodb.PutItemInput{
		Item:      dynamoMap,
		TableName: aws.String(tableName),
	}

	res, err := dynamoClient.PutItem(creationRequest)
	if err != nil && res != nil {
		return nil, errors.New("Could not create this item")
	}

	return &newUser, nil
}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynamoClient dynamodbiface.DynamoDBAPI) (
	*User, error,
) {
	var u User
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New("Invalid email")
	}

	currentUser, _ := GetUser(u.Email, tableName, dynamoClient)
	if currentUser != nil {
		return nil, errors.New("User does not exist")
	}

	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, errors.New("Could not convert item")
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynamoClient.PutItem(input)
	if err != nil {
		return nil, errors.New("Could not update item")
	}

	return &u, nil
}

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynamoClient dynamodbiface.DynamoDBAPI) error {
	email := req.QueryStringParameters["email"]
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := dynamoClient.DeleteItem(input)

	if err != nil {
		return errors.New("Could not delete item")
	}

	return nil
}
