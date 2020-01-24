package validation

import (
	"errors"
	"reflect"
	"testing"

	"github.com/illuminateeducation/rest-service-lib-go/pkg/types"
	"gopkg.in/go-playground/validator.v9"
)

func TestTransformErrors(t *testing.T) {
	pgv := PlaygroundValidator{}

	t.Run("Check for edge case error", func(t *testing.T) {
		err := &validator.InvalidValidationError{}
		actualErr := pgv.transformErrors(err)

		if err != actualErr {
			t.Error("transformErrors should always return InvalidValidationError when they occur")
		}
	})

	t.Run("If there aren't any errors, return nil", func(t *testing.T) {
		err := validator.ValidationErrors{}
		if pgv.transformErrors(err) != nil {
			t.Error("transformErrors should return nil if the list of errors is empty")
		}
	})

	t.Run("Error values should be consolidated to a single error", func(t *testing.T) {
		errs := validator.ValidationErrors{MockFieldError{}, MockFieldError{}}
		actualErr := pgv.transformErrors(errs)

		if _, ok := actualErr.(validator.ValidationErrors); ok {
			t.Error("The returned value was not an error interface")
		}
	})
}

func TestNewPlaygroundValidator(t *testing.T) {
	value := NewPlaygroundValidator()

	_, ok := value.(Validator)
	if !ok {
		t.Error("NewPlaygroundValidator isn't returning a Validator type")
	}
}

func TestPlaygroundValidator_Var(t *testing.T) {
	expectedError := errors.New("test error")

	pv := PlaygroundValidator{
		validator: MockPlaygroundValidatorAdaptor{
			VarReturnVal: expectedError,
		},
	}

	t.Run("PlaygroundValidator::Var will return a transformed error", func(t *testing.T) {
		if pv.Var("value", "-") != expectedError {
			t.Error("The expected error wasn't returned from PlaygroundValidator::Var")
		}
	})
}

func TestPlaygroundValidator_Struct(t *testing.T) {
	expectedError := errors.New("test error")

	pv := PlaygroundValidator{
		validator: MockPlaygroundValidatorAdaptor{
			StructReturnVal: expectedError,
		},
	}

	t.Run("PlaygroundValidator::Var will return a transformed error", func(t *testing.T) {
		if pv.Struct("value") != expectedError {
			t.Error("The expected error wasn't returned from PlaygroundValidator::Struct")
		}
	})
}

func TestPlaygroundUuidNotBlank(t *testing.T) {
	t.Run("A valid UUIDv4 should pass", func(t *testing.T) {
		m := MockFieldLevel{
			FieldValue: "f99798eb-6724-4448-9ee0-9259cde51c2a",
		}

		if !playgroundUuidNotBlank(m) {
			t.Errorf("A valid UUIDv4 should return true but got false")
		}
	})
}

func TestPlaygroundStringNotBlank(t *testing.T) {

	tests := []struct {
		name  string
		param string
	}{
		{name: "Too few parameters in notblank", param: ""},
		{name: "Too many parameters in notblank", param: "1 2 3"},
		{name: "first parameter panics on non-integer", param: "bad -"},
		{name: "first parameter panics on negative integer", param: "-1 -"},
		{name: "second parameter panics on non-integer", param: "- bad"},
		{name: "second parameter panics on negative integer", param: "- -1"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Invalid notblank settings. No panic was raised.")
				}
			}()

			m := MockFieldLevel{
				FieldValue:  "test",
				StructParam: test.param,
			}
			playgroundStringNotBlank(m)
		})
	}

	t.Run("Default parameters notblank, is a valid check", func(t *testing.T) {
		m := MockFieldLevel{
			FieldValue:  "test",
			StructParam: "1 10",
		}
		if !playgroundStringNotBlank(m) {
			t.Error("A valid non-blank string returned false")
		}
	})

}

func TestValidateValuer(t *testing.T) {
	t.Run("Return nil if its not a driver.Valuer", func(t *testing.T) {
		if ValidateValuer(reflect.ValueOf("not a valuer")) != nil {
			t.Error("A non-valuer should return nil")
		}
	})

	t.Run("Return value of the driver.Valuer", func(t *testing.T) {
		expectedString := "test"
		s := types.NullString{}
		s.Scan(expectedString)
		actual := ValidateValuer(reflect.ValueOf(s))
		if actual != expectedString {
			t.Errorf("Didn't return expected value: %v", actual)
		}
	})

	t.Run("Return a zero if the driver.Valuer returns nil.", func(t *testing.T) {
		expectedReturn := 0
		actual := ValidateValuer(reflect.ValueOf(types.NullString{}))
		if actual != expectedReturn {
			t.Errorf("Didn't return expected value: %v", actual)
		}
	})
}

func TestJSONTagNameFunc(t *testing.T) {

	type NoJSONTag struct {
		Name string
	}

	type LowercaseJSONTag struct {
		Name string `json:"name"`
	}

	type DifferentJSONTag struct {
		Name string `json:"not_name"`
	}

	type SkippedJSONTag struct {
		Name string `json:"-"`
	}

	tests := []struct {
		name         string
		testCase     interface{}
		expectedName string
	}{
		{name: "Struct with no json tag", testCase: NoJSONTag{}, expectedName: "Name"},
		{name: "Struct with skipped(-) json tag", testCase: SkippedJSONTag{}, expectedName: "Name"},
		{name: "Struct with lower case json tag", testCase: LowercaseJSONTag{}, expectedName: "name"},
		{name: "Struct with different json tag", testCase: DifferentJSONTag{}, expectedName: "not_name"},
	}

	for _, test := range tests {
		st := reflect.TypeOf(test.testCase)
		name := JSONTagNameFunc(st.Field(0))

		if name != test.expectedName {
			t.Errorf(`%s: Expected "%s", got "%s"`, test.name, test.expectedName, name)
		}
	}
}
