package setup

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func NewDynamoClient(s *session.Session) *dynamodb.DynamoDB {
	client := dynamodb.New(s)
	return client
}
