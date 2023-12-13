package setup

import (
	"github.com/diontr00/serverlessgo/translator"
	"github.com/diontr00/serverlessgo/validator"
)

// set up new validator
func NewValidator(trans translator.Translator) validator.Validator {
	return validator.New(trans)
}
