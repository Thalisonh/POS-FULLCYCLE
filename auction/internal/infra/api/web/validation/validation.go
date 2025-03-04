package validation

import (
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/validator/v10"
	"github.com/thalisonh/auction/configuration/rest_err"

	ut "github.com/go-playground/universal-translator"
	validator_en "github.com/go-playground/validator/v10/translations/en"
)

var (
	Validate = validator.New()
	transl   ut.Translator
)

func init() {
	if value, ok := binding.Validator.Engine().(*validator.Validate); ok {
		en := en.New()
		enTransl := ut.New(en, en)
		transl, _ = enTransl.GetTranslator("en")
		validator_en.RegisterDefaultTranslations(value, transl)
	}
}

func ValidateErr(validationErr error) *rest_err.RestErr {
	var jsonErr *json.UnmarshalTypeError
	var jsonValidation validator.ValidationErrors

	if errors.As(validationErr, &jsonErr) {
		return rest_err.NewBadRequestError("Invalid type error")
	} else if errors.As(validationErr, &jsonValidation) {
		errorCauses := []rest_err.Causes{}

		for _, e := range validationErr.(validator.ValidationErrors) {
			errorCauses = append(errorCauses, rest_err.Causes{
				Field:   e.Field(),
				Message: e.Translate(transl),
			})

			return rest_err.NewBadRequestError("Invalid input", errorCauses...)
		}
	} else {
		return rest_err.NewBadRequestError("Error trying to validate input")
	}

	return nil
}
