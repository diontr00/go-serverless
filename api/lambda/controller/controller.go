package Controller

import (
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/diontr00/serverlessgo/model"
	"github.com/rs/zerolog"

	"github.com/goccy/go-json"
)

type jsonBindError struct {
	Field  any `json:"field"`
	Got    any `json:"got"`
	Expect any `json:"expect"`
}

type Controller struct {
	l zerolog.Logger
	u model.UserUseCase
}

func (c *Controller) Get(req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	email := req.QueryStringParameters["email"]
	lang := req.QueryStringParameters["locale"]
	return c.u.GetUser(lang, email)
}

func (c *Controller) Post(req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	lang := req.QueryStringParameters["locale"]
	var user model.User
	errResp := c.Bind(lang, req, &user)
	if errResp != nil {
		return errResp, nil
	}
	return c.u.CreateUser(lang, user)
}

func (c *Controller) Put(req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	lang := req.QueryStringParameters["locale"]

	var user model.User
	errResp := c.Bind(lang, req, &user)
	if errResp != nil {
		return errResp, nil
	}
	return c.u.UpdateUser(lang, user)

}

func (c *Controller) Delete(req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	email := req.QueryStringParameters["email"]
	lang := req.QueryStringParameters["locale"]
	return c.u.DeleteUser(lang, email)

}

func (c *Controller) Unhandler(req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	lang := req.QueryStringParameters["locale"]
	return c.u.UnHandler(lang)
}

// Bind the request payload to struct i , if any error occur response will be return else nil
func (c *Controller) Bind(lang string, req *events.APIGatewayProxyRequest, i interface{}) *events.APIGatewayProxyResponse {

	body := strings.NewReader(req.Body)

	err := json.NewDecoder(body).Decode(i)
	if err != nil {
		var res *events.APIGatewayProxyResponse

		switch e := err.(type) {
		case *json.UnmarshalTypeError:
			res, _ = c.u.Json(lang, http.StatusBadRequest, map[string]jsonBindError{"error": {Field: e.Field, Expect: e.Type.Name(), Got: e.Value}})
		case *json.SyntaxError:
			res, _ = c.u.Json(lang, http.StatusBadRequest, map[string]string{"error": "supplied json syntax is invalid"})
		default:
			c.l.Err(e).Msg("Unknown error when decode the body")
			res, _ = c.u.Json(lang, http.StatusBadRequest, map[string]string{"error": "internal server error"})
		}
		return res
	}
	return nil

}

func NewController(u model.UserUseCase, l zerolog.Logger) *Controller {
	return &Controller{u: u, l: l}
}
