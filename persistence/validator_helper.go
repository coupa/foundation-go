package persistence

import (
	"github.com/asaskevich/govalidator"
	"reflect"
)

func Validate(obj interface{}) (ValidationErrors, error) {
	ret := ValidationErrors{Errors: map[string][]string{}}
	_, err := govalidator.ValidateStruct(obj)
	if err != nil {
		switch err.(type) {
		case govalidator.Error:
			addValidationError(err.(govalidator.Error), &ret)
		case govalidator.Errors:
			for _, err := range err.(govalidator.Errors) {
				switch err.(type) {
				case govalidator.Error:
					addValidationError(err.(govalidator.Error), &ret)
				default:
					return ret, err
				}
			}
		}
		if reflect.TypeOf(err) == reflect.TypeOf(govalidator.Errors{}) {
			govalidatorErrors := err.(govalidator.Errors)
			govalidatorErrors.Errors()
		}
	}
	return ret, nil
}

func addValidationError(govalidatorError govalidator.Error, validationErrors *ValidationErrors) {
	var slice []string
	var exists bool
	if slice, exists = validationErrors.Errors[govalidatorError.Name]; !exists {
		slice = []string{}
	}
	validationErrors.Errors[govalidatorError.Name] = append(slice, govalidatorError.Err.Error())
}
