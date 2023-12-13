package repository

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/diontr00/serverlessgo/config/env"
	"github.com/diontr00/serverlessgo/model"
)

type userRepo struct {
	env        env.DynamoEnv
	dynaClient dynamodbiface.DynamoDBAPI
}

func (r *userRepo) GetUser(email string) (*model.User, error) {

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: &r.env.TableName,
	}

	result, err := r.dynaClient.GetItem(input)

	if err != nil {
		return nil, fmt.Errorf("Error fetching record: %v", err)
	}

	item := new(model.User)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling: %v", err)
	}

	return item, nil
}

func (r *userRepo) GetUsers() (*[]model.User, error) {
	input := &dynamodb.ScanInput{
		TableName: &r.env.TableName,
	}

	res, err := r.dynaClient.Scan(input)
	if err != nil {
		return nil, fmt.Errorf("Error fetching all : %v", err)
	}
	users := new([]model.User)
	err = dynamodbattribute.UnmarshalListOfMaps(res.Items, users)
	if err != nil {

		return nil, fmt.Errorf("Error marshalling: %v", err)
	}

	return users, nil
}

func (r *userRepo) CreateUser(user model.User) error {
	current, _ := r.GetUser(user.Email)
	if current != nil {
		return fmt.Errorf("Error already exist")
	}

	m, err := dynamodbattribute.MarshalMap(user)

	if err != nil {
		return fmt.Errorf("Error cannot marshalling item")
	}

	input := &dynamodb.PutItemInput{
		Item:      m,
		TableName: &r.env.TableName,
	}

	_, err = r.dynaClient.PutItem(input)
	if err != nil {
		return fmt.Errorf("Error create user")
	}
	return nil
}

func (r *userRepo) UpdateUser(user model.User) error {

	current, _ := r.GetUser(user.Email)
	if current == nil {
		return fmt.Errorf("Error not exist")
	}

	m, err := dynamodbattribute.MarshalMap(user)

	if err != nil {
		return fmt.Errorf("Error cannot marshalling item")
	}

	input := &dynamodb.PutItemInput{
		Item:      m,
		TableName: &r.env.TableName,
	}

	_, err = r.dynaClient.PutItem(input)
	if err != nil {
		return fmt.Errorf("Error create user")
	}
	return nil

}

func (r *userRepo) DeleteUser(email string) error {
	current, _ := r.GetUser(email)
	if current == nil {
		return fmt.Errorf("Error not exist")
	}

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},

		TableName: &r.env.TableName,
	}

	_, err := r.dynaClient.DeleteItem(input)
	if err != nil {
		return fmt.Errorf("Error delete user")
	}
	return nil
}

func NewRepository(e env.DynamoEnv, client dynamodbiface.DynamoDBAPI) *userRepo {
	return &userRepo{
		env:        e,
		dynaClient: client,
	}
}
