package validation

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// StringNotBlank will check that the value is either null or the length is greater than zero if a 0 was passed in
// then the value was null
func StringNotBlank(field reflect.Value, min *int, max *int) bool {
	if field.Kind() == reflect.Int {
		return field.Interface() == 0
	}

	if field.Type().String() == "string" {
		val, _ := field.Interface().(string)
		if len(val) == 0 {
			return false
		}

		if min != nil && len(val) < *min {
			return false
		}

		if max != nil && len(val) > *max {
			return false
		}
	}

	return true
}

// UuidNotBlank will check that the value is either null or a UUIDv4.
// If a "0" was passed in then the value was null.
func UuidNotBlank(field reflect.Value) bool {
	if field.Kind() == reflect.Int {
		return field.Interface() == 0
	}

	if field.Type().String() == "string" {
		val := field.Interface().(string)
		validate := validator.New()
		err := validate.Var(val, "uuid4")

		if err != nil {
			return false
		}
	}

	return true
}
