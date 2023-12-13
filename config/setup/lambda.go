package setup

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func Newlambda(s *session.Session) *lambda.Lambda {
	return lambda.New(s)
}
