package model

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

type JsonFieldExpect struct {
	Field  any `json:"field"`
	Expect any `json:"expect"`
	Got    any `json:"got"`
}

// Error to be return when json binding fail
type JsonBindingError struct {
	Err JsonFieldExpect `json:"error"`
}

func (j JsonBindingError) Error() string {
	return fmt.Sprintf("%v", j.Err)
}

// Error when validating the request payload
type FieldError interface {
	// Return the validation tag that failed for example "required or lt"
	Tag() string
	// Returns the field name that failed validation
	Field() string
	// Returns the param that send with the request in case needed for creating the message
	Param() string
}

// usecase will perform the main logic that will be invovked in the controller layout
type UserUseCase interface {
	GetUser(lang, email string) (*events.APIGatewayProxyResponse, error)
	GetUsers(lang string) (*events.APIGatewayProxyResponse, error)
	CreateUser(lang string, user User) (*events.APIGatewayProxyResponse, error)

	UpdateUser(lang string, user User) (*events.APIGatewayProxyResponse, error)

	DeleteUser(lang, email string) (*events.APIGatewayProxyResponse, error)

	UnHandler(lang string) (*events.APIGatewayProxyResponse, error)
	Json(lang string, status int, body interface{}) (*events.APIGatewayProxyResponse, error)
}
type UserRepo interface {
	GetUser(email string) (*User, error)
	GetUsers() (*[]User, error)
	CreateUser(user User) error
	UpdateUser(user User) error
	DeleteUser(email string) error
}

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type FieldErrorResponse struct {
	Field string `json:"field"`
	Msg   string `json:"error_msg"`
}

type UserResponse struct {
	Error error `json:"error,omitempty"`
	Data  any   `json:"data,omitempty"`
}
