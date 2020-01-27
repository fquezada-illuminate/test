package validation

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/illuminateeducation/rest-service-lib-go/pkg/types"
	"gopkg.in/go-playground/validator.v9"
)

type PlaygroundValidatorAdaptor interface {
	Var(field interface{}, tag string) error
	Struct(s interface{}) error
}

type PlaygroundValidator struct {
	validator PlaygroundValidatorAdaptor
}

func (pgv *PlaygroundValidator) Var(field interface{}, options interface{}) error {
	err := pgv.validator.Var(field, options.(string))

	return pgv.transformErrors(err)
}

func (pgv *PlaygroundValidator) Struct(s interface{}) error {
	// err := pgv.validator.Struct(s)

	return nil
}

// TransformErrors will take the errors from the playground library into a standard error
func (pgv PlaygroundValidator) transformErrors(err error) error {

	// this check is only needed when your code could produce
	// an invalid value for validation such as interface with nil
	// value most including myself do not usually have code like this.
	if _, ok := err.(*validator.InvalidValidationError); ok {
		return err
	}

	// if a regular error is passed in, return the error
	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	errs := make([]Error, 0)
	for _, e := range validationErrs {
		errs = append(errs, newFieldError(e.Field(), e.ActualTag(), e.Value()))
	}

	if len(errs) == 0 {
		return nil
	}

	return consolidateValidationErrors(errs)
}

// NewPlaygroundValidator instantiates a new instance of the PlaygroundValidator.  It is responsible for adding
// any additional functions and tags to the base library
func NewPlaygroundValidator() Validator {
	v := validator.New()
	v.RegisterValidation("notblank", playgroundStringNotBlank)
	v.RegisterValidation("uuidnotblank", playgroundUuidNotBlank)
	v.RegisterCustomTypeFunc(ValidateValuer, types.NullString{})
	v.RegisterTagNameFunc(JSONTagNameFunc)

	return &PlaygroundValidator{
		validator: v,
	}
}

func playgroundUuidNotBlank(fl validator.FieldLevel) bool {
	return UuidNotBlank(fl.Field())
}

// playgroundStringNotBlank is a callback specifically for the PlaygroundValidator
func playgroundStringNotBlank(fl validator.FieldLevel) bool {
	params := strings.Split(fl.Param(), " ")

	if len(params) > 2 {
		panic("too many parameters for notblank validation tag")
	}

	if len(params) < 2 {
		panic("notblank validation tag requires 2 parameters")
	}

	var min *int
	if params[0] != "-" {
		m, err := strconv.Atoi(params[0])
		if err != nil {
			panic(fmt.Sprintf("notblank parameter was not an integer, got: %s", params[0]))
		}

		if m < 0 {
			panic("notblank min parameter cannot be negative")
		}

		min = &m
	}

	var max *int
	if params[1] != "-" {
		m, err := strconv.Atoi(params[1])
		if err != nil {
			panic(fmt.Sprintf("notblank parameter was not an integer, got: %s", params[1]))
		}
		if m < 0 {
			panic("notblank max parameter cannot be negative")
		}

		max = &m
	}

	return StringNotBlank(fl.Field(), min, max)
}

// ValidateValuer implements validator.CustomTypeFunc
func ValidateValuer(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(driver.Valuer); ok {
		val, err := valuer.Value()
		if err == nil {
			if val == nil {
				return 0
			}
			return val
		}
	}

	return nil
}

// JSONTagNameFunc uses the json tag for naming in errors.
// The field name is returned if no json tag is present.
// An empty string is returned if the field is to be skipped eg. `json:"-"``
// from example: https://godoc.org/gopkg.in/go-playground/validator.v9#Validate.RegisterTagNameFunc
func JSONTagNameFunc(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

	// if there is no defined json tag or the tag is to be skipped fall back to the field name
	if name == "" || name == "-" {
		return fld.Name
	}

	return name
}
