package config

import (
	"context"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"io"
	oslog "log"
	"os"
	"os/signal"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/diontr00/serverlessgo/api/lambda/router"
	"github.com/diontr00/serverlessgo/config/env"
	. "github.com/diontr00/serverlessgo/config/setup"
	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
)

type Applications struct {
	Lambda  *router.LamdaRouter
	Dynamo  *dynamodb.DynamoDB
	Logger  *zerolog.Logger
	Env     *env.Env
	Logfile io.WriteCloser
}

// Start rest server and register clean up function
func (a *Applications) Start() {
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)

	go func() {
		<-terminate
		oslog.Println("Gratefully Shutdown , Doing Cleanup Task...😷")

		ctx, cancel := context.WithTimeout(context.Background(), a.Env.App.ShutdownTimeout)
		defer func() {
			cancel()
			a.ShutDown(ctx)

		}()
	}()

	lambda.Start(a.Lambda.Handler)

}

func (a *Applications) ShutDown(ctx context.Context) {
	select {
	case <-ctx.Done():
		log.Fatal().Msgf("Error Shutdown %v", ctx.Err())
	default:

		if a.Logfile != nil {
			err := a.Logfile.Close()
			if err != nil {

				log.Fatal().Msgf("Closing log file error : %v", ctx.Err())
			}

		}

	}

}

func NewApp(ctx context.Context) *Applications {

	var app = new(Applications)
	select {
	case <-ctx.Done():
		log.Fatal().Msgf("Setup Application timeout : %v", ctx.Err())

	default:

		env := env.NewEnv(ctx)
		translator := NewTranslator()
		validator := NewValidator(translator)
		logger, logfile := NewLogger(env)
		session := NewAwsSession(&env.Aws)
		dynamo := NewDynamoClient(session)

		router := router.NewLambdaRouter(logger, dynamo, env, translator, validator)

		app.Env = env
		app.Lambda = router
		app.Logger = &logger
		app.Logfile = logfile
		app.Dynamo = dynamo

	}

	return app
}
