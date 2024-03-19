package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	customValidator *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	return &CustomValidator{
		customValidator: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (v *CustomValidator) Validate(s interface{}) *validator.ValidationErrors{
	err, ok := v.customValidator.Struct(s).(validator.ValidationErrors)
	if !ok {
		return nil
	}
	return &err
}

func (v *CustomValidator) MakePrettyErrors(err *validator.ValidationErrors) string{
	validationErrors := ""

	for _, err := range *err {
		switch err.Tag() {
		case "required":
			validationErrors = validationErrors + fmt.Sprintf("The field '%s' is required;", err.Field())
		case "email":
			validationErrors = validationErrors + fmt.Sprintf("Invalid email format in field '%s';", err.Field())
		case "min":
			validationErrors = validationErrors + fmt.Sprintf("The field '%s' must be at least %s characters long;", err.Field(), err.Param())
		}
	}
	return validationErrors
}