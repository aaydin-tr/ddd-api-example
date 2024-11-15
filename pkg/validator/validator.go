package validator

import (
	"github.com/aaydin-tr/gowit-case/interface/http/response"
	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func New() *CustomValidator {
	return &CustomValidator{Validator: validator.New()}
}

func MsgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "uuid4":
		return "This field must be valid uuid4"
	}
	return fe.Error()
}

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.Validator.Struct(i)
	var errors []*response.ValidationMessage
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element response.ValidationMessage
			element.FailedField = err.Field()
			element.Tag = err.Tag()
			element.Message = MsgForTag(err)
			errors = append(errors, &element)
		}
	}

	if len(errors) > 0 {
		return &response.ErrorResponse{
			Message: "Validation error",
			Errors:  errors,
			Status:  400,
		}
	}
	return err
}
