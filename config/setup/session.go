package setup

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/diontr00/serverlessgo/config/env"
)

func NewAwsSession(env *env.AwsEnv) *session.Session {
	awsSession, err := session.NewSession(&aws.Config{
		Region: &env.Region,
	})
	if err != nil {
		log.Fatalf("[Error] Creating Session File %v \n", err)
	}
	return awsSession
}
