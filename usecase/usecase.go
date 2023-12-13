package usecase

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/diontr00/serverlessgo/model"
	"github.com/diontr00/serverlessgo/translator"
	"github.com/diontr00/serverlessgo/validator"
	"github.com/rs/zerolog"

	"github.com/goccy/go-json"
	"github.com/hashicorp/errwrap"
)

type userUsecase struct {
	l zerolog.Logger
	v validator.Validator
	t translator.Translator
	r model.UserRepo
}

func (u *userUsecase) GetUser(lang string, email string) (*events.APIGatewayProxyResponse, error) {

	errs := u.v.ValidateAndTranslateAny(lang, email, "required,email")
	if errs != nil {
		return u.Json(lang, http.StatusBadRequest, errs)
	}

	res, err := u.r.GetUser(email)

	if err != nil {

		if errwrap.Contains(err, "fetching") {
			return u.Json(lang, http.StatusBadRequest, model.UserResponse{Error: err})
		}

	}

	return u.Json(lang, http.StatusOK, model.UserResponse{Data: res})

}

func (u *userUsecase) GetUsers(lang string) (*events.APIGatewayProxyResponse, error) {

	res, err := u.r.GetUsers()

	if err != nil {
		return u.Json(lang, http.StatusBadRequest, model.UserResponse{Error: err})
	}

	return u.Json(lang, http.StatusOK, model.UserResponse{Data: res})

}

func (u *userUsecase) CreateUser(lang string, user model.User) (*events.APIGatewayProxyResponse, error) {
	errs := u.v.ValidateRequestAndTranslate(lang, &user)

	if errs != nil {
		return u.Json(lang, http.StatusBadRequest, errs)
	}

	err := u.r.CreateUser(user)

	if err != nil {
		return u.Json(lang, http.StatusBadRequest, model.UserResponse{Error: err})
	}

	return u.Json(lang, http.StatusOK, model.UserResponse{Data: fmt.Sprintf("Created %s", user.Email)})
}

func (u *userUsecase) UpdateUser(lang string, user model.User) (*events.APIGatewayProxyResponse, error) {

	errs := u.v.ValidateRequestAndTranslate(lang, &user)
	if errs != nil {
		return u.Json(lang, http.StatusBadRequest, errs)
	}

	err := u.r.UpdateUser(user)

	if err != nil {
		return u.Json(lang, http.StatusBadRequest, model.UserResponse{Error: err})
	}

	return u.Json(lang, http.StatusOK, model.UserResponse{Data: fmt.Sprintf("Updated %s", user.Email)})

}

func (u *userUsecase) DeleteUser(lang, email string) (*events.APIGatewayProxyResponse, error) {

	errs := u.v.ValidateAndTranslateAny(lang, email, "required,email")
	if errs != nil {
		return u.Json(lang, http.StatusBadRequest, errs)
	}

	err := u.r.DeleteUser(email)

	if err != nil {
		return u.Json(lang, http.StatusBadRequest, model.UserResponse{Error: err})
	}

	return u.Json(lang, http.StatusOK, model.UserResponse{Data: fmt.Sprintf("Deleted %s", email)})

}

func (u *userUsecase) UnHandler(lang string) (*events.APIGatewayProxyResponse, error) {

	return u.Json(lang, http.StatusMethodNotAllowed, "method not allowed")
}

// Serialize body into api response
func (u *userUsecase) Json(lang string, status int, body interface{}) (*events.APIGatewayProxyResponse, error) {

	resp := events.APIGatewayProxyResponse{Headers: map[string]string{
		"Content-type": "application/json",
	}}
	resp.StatusCode = status
	byteBody, err := json.Marshal(body)
	if err != nil {
		u.l.Err(err).Msgf("Unknown error when try to send the body %v", body)
		resp.StatusCode = http.StatusInternalServerError
		msg, err := u.t.TranslateMessage(lang, "internalError", nil, nil)
		if err != nil {
			u.l.Err(err).Msgf("Cloud not trranslate message internalError for lang %s", lang)
			msg = "internal server error"
		}

		resp.Body = msg
		return &resp, err
	}
	resp.Body = string(byteBody)
	return &resp, nil
}

func NewUseruc(l zerolog.Logger, v validator.Validator, t translator.Translator, r model.UserRepo) model.UserUseCase {
	return &userUsecase{l: l, v: v, t: t, r: r}
}
