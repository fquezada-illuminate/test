package validation

import (
	"github.com/go-playground/universal-translator"
	"reflect"
)

type MockValidationError struct {
	field string
	rule  string
	value string
	param string
}

func (m MockValidationError) Field() string {
	return m.field
}
func (m MockValidationError) Rule() string {
	return m.rule
}
func (m MockValidationError) Value() interface{} {
	return m.value
}
func (m MockValidationError) Param() interface{} {
	return m.param
}

type MockFieldError struct {
}

func (MockFieldError) ActualTag() string {
	return "tag"
}

func (MockFieldError) Field() string {
	return "field-name"
}

func (MockFieldError) Kind() reflect.Kind {
	panic("implement me")
}

func (MockFieldError) Namespace() string {
	panic("implement me")
}

func (m MockFieldError) Param() string {
	return m.Param()
}

func (MockFieldError) StructField() string {
	panic("implement me")
}

func (MockFieldError) StructNamespace() string {
	panic("implement me")
}

func (MockFieldError) Tag() string {
	panic("implement me")
}

func (MockFieldError) Translate(ut ut.Translator) string {
	panic("implement me")
}

func (MockFieldError) Type() reflect.Type {
	panic("implement me")
}

func (MockFieldError) Value() interface{} {
	return "value"
}

type MockPlaygroundValidatorAdaptor struct {
	VarReturnVal    error
	StructReturnVal error
}

func (m MockPlaygroundValidatorAdaptor) Var(field interface{}, tag string) error {
	return m.VarReturnVal
}

func (m MockPlaygroundValidatorAdaptor) Struct(s interface{}) error {
	return m.StructReturnVal
}

type MockFieldLevel struct {
	FieldValue  interface{}
	StructParam string
}

func (m MockFieldLevel) Top() reflect.Value {
	panic("implement me")
}

func (m MockFieldLevel) Parent() reflect.Value {
	panic("implement me")
}

func (m MockFieldLevel) Field() reflect.Value {
	return reflect.ValueOf(m.FieldValue)
}

func (m MockFieldLevel) FieldName() string {
	panic("implement me")
}

func (m MockFieldLevel) StructFieldName() string {
	panic("implement me")
}

func (m MockFieldLevel) Param() string {
	return m.StructParam
}

func (m MockFieldLevel) ExtractType(field reflect.Value) (value reflect.Value, kind reflect.Kind, nullable bool) {
	panic("implement me")
}

func (m MockFieldLevel) GetStructFieldOK() (reflect.Value, reflect.Kind, bool) {
	panic("implement me")
}
