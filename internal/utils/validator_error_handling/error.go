package validator_error_handling

import "github.com/go-playground/validator/v10"

var validationErrors = map[string]string{
	"required":   "поле обязательно для заполнения",
	"email":      "некорректный email",
	"min":        "значение должно быть больше минимального",
	"max":        "значение должно быть меньше максимального",
	"uuid":       "некорректный формат UUID",
	"oneof":      "недопустимое значение",
	"unique":     "значение должно быть уникальным",
	"exists":     "значение должно существовать",
	"not_exists": "значение не должно существовать",
	"not_empty":  "значение не должно быть пустым",
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrorList struct {
	errors []ValidationError
}

func newValidationErrorList(count int) ValidationErrorList {
	return ValidationErrorList{
		errors: make([]ValidationError, 0, count),
	}
}

func (v *ValidationErrorList) AddError(field string, tag string) {
	errorMessage, ok := validationErrors[tag]
	if !ok {
		errorMessage = "неизвестная ошибка валидации"
	}

	v.errors = append(v.errors, ValidationError{
		Field:   field,
		Message: errorMessage,
	})
}

func (v *ValidationErrorList) Errors() []ValidationError {
	return v.errors
}

func (v *ValidationErrorList) HasErrors() bool {
	return len(v.errors) > 0
}

func ValidateError(err error) *ValidationErrorList {
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		validationErrorList := newValidationErrorList(len(validationErrs))

		for _, e := range validationErrs {
			validationErrorList.AddError(e.Field(), e.Tag())
		}

		return &validationErrorList
	}

	return nil
}
