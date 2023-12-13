package router

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/diontr00/serverlessgo/api/lambda/controller"
	"github.com/diontr00/serverlessgo/config/env"
	"github.com/diontr00/serverlessgo/repository"
	"github.com/diontr00/serverlessgo/usecase"
	"github.com/rs/zerolog"

	"github.com/diontr00/serverlessgo/translator"
	"github.com/diontr00/serverlessgo/validator"
)

type LamdaRouter struct {
	Translator translator.Translator
	Validator  validator.Validator
	controller *Controller.Controller
}

func (r *LamdaRouter) Handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	locale := req.QueryStringParameters["locale"]

	switch locale {
	case "en":
		locale = "en-US"
	case "vi":
		locale = "vi-VN"
	default:
		locale = "en-US"
	}
	req.PathParameters[locale] = locale

	switch req.HTTPMethod {
	case "GET":
		return r.controller.Get(&req)

	case "POST":
		return r.controller.Post(&req)
	case "PUT":
		return r.controller.Put(&req)
	case "DELETE":
		return r.controller.Delete(&req)
	default:
		return r.controller.Unhandler(&req)
	}
}

func NewLambdaRouter(logger zerolog.Logger, d dynamodbiface.DynamoDBAPI, e *env.Env, t translator.Translator, v validator.Validator) *LamdaRouter {

	userRepo := repository.NewRepository(e.Dynamo, d)
	userUsecase := usecase.NewUseruc(logger, v, t, userRepo)
	userController := Controller.NewController(userUsecase, logger)

	return &LamdaRouter{
		Translator: t,
		Validator:  v,
		controller: userController,
	}
}
