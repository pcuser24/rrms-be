package validation

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type XValidator struct {
	*validator.Validate
}

type ErrorResponse struct {
	Error       bool
	FailedField string
	Tag         string
	Value       interface{}
}

type ErrorMessage struct {
	Message string `json:"message"`
}

var defaultStructValidator = validator.New()

func GetDefaultValidator() *validator.Validate {
	copy := *defaultStructValidator
	return &copy
}

func init() {
	log.Println("Setting up UUID validator...")
	// validate UUID
	defaultStructValidator.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
		if valuer, ok := field.Interface().(uuid.UUID); ok {
			return valuer.String()
		}
		return nil
	}, uuid.Nil)
}

func ValidateStruct(v *validator.Validate, data interface{}) []ErrorResponse {
	if v == nil {
		v = defaultStructValidator
	}
	validationErrors := []ErrorResponse{}
	errs := v.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			var elem ErrorResponse

			elem.FailedField = err.Field() // Export struct field name
			elem.Tag = err.Tag()           // Export struct tag
			elem.Value = err.Value()       // Export field value
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}

func GetValidationError(errs []ErrorResponse) string {
	errMsgs := make([]string, 0)

	for _, err := range errs {
		errMsgs = append(errMsgs, fmt.Sprintf(
			"[%s]: '%v' | Needs to implement '%s'",
			err.FailedField,
			err.Value,
			err.Tag,
		))
	}

	return strings.Join(errMsgs, "; ")
}
