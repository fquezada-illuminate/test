package validation

import (
	"errors"
	"fmt"
	"strings"
)

const ERROR_DELIMITER = " || "

type Error interface {
	Field() string
	Rule() string
	Value() interface{}
}

type FieldError struct {
	field string
	rule  string
	value interface{}
}

func (fe FieldError) Rule() string {
	return fe.rule
}

func (fe FieldError) Value() interface{} {
	return fe.value
}

func (fe FieldError) Field() string {
	return fe.field
}

func newFieldError(field string, rule string, value interface{}) FieldError {
	return FieldError{
		field: field,
		rule:  rule,
		value: value,
	}
}

// ConsolidateValidationErrors will take all of the accumulated errors and combine them as a single error
// separated by the delimiter
func consolidateValidationErrors(errs []Error) error {

	errStrings := make([]string, 0)
	for _, e := range errs {
		errStrings = append(errStrings, fmt.Sprintf("%s: Failed validation for %s with value %v", e.Field(), e.Rule(), e.Value()))
	}

	return errors.New(strings.Join(errStrings, ERROR_DELIMITER))
}
