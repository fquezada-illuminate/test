package validation

import (
	"encoding/json"
	"testing"
)

func TestDecodeRequest(t *testing.T) {
	type TestStruct struct {
		Name     string `json:"name"`
		DoNotSet string `json:"doNotSet"`
	}

	validationMap := []string{
		"Name",
	}

	testStruct := TestStruct{}

	t.Run("The obj should be a pointer", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fail()
			}
		}()
		notPtr := "not a pointer"
		DecodeRequest([]byte(""), validationMap, notPtr)
	})

	t.Run("Return an error if cant unmarshal body into map[string]json.RawMessage", func(t *testing.T) {
		err := DecodeRequest([]byte("test"), validationMap, &testStruct)
		if err == nil {
			t.Fail()
		}
	})

	t.Run("The validFields must contain values that are part of the struct ", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fail()
			}
		}()
		DecodeRequest([]byte(`{"name":"value"}`), []string{"Invalid"}, &testStruct)
	})

	t.Run("A key in the response isn't part of the struct based on the json tag", func(t *testing.T) {
		err := DecodeRequest([]byte(`{"invalid":"value"}`), []string{}, &testStruct)
		if err == nil {
			t.Fail()
		}
	})

	t.Run("A key that cannot be set is passed in the request", func(t *testing.T) {
		err := DecodeRequest([]byte(`{"doNotSet":"value"}`), validationMap, &testStruct)
		if err == nil {
			t.Fail()
		}
	})

	t.Run("Catch a list of unmarshalling errors for the individual fields", func(t *testing.T) {
		err := DecodeRequest([]byte(`{"name": 123}`), []string{}, &testStruct)
		if err == nil {
			t.Fail()
		}
	})

	t.Run("Properly decode the struct members from the request", func(t *testing.T) {
		DecodeRequest([]byte(`{"name":"value"}`), []string{}, &testStruct)
		if testStruct.Name != "value" {
			t.Errorf("The struct fields were not properly decoded: %s", testStruct.Name)
		}
	})

	t.Run("Json NUll values should default to the empty value of the struct", func(t *testing.T) {
		DecodeRequest([]byte(`{"name":null}`), []string{}, &testStruct)
		if testStruct.Name != "" {
			t.Errorf("The struct fields were not properly decoded: %s", testStruct.Name)
		}
	})

	t.Run("Properly decode a complex request struct from the request", func(t *testing.T) {
		type nestStruct struct {
			Test string `json:"test"`
		}

		type request struct {
			Name         string     `json:"name"`
			Age          int        `json:"age"`
			Active       bool       `json:"active"`
			NestedStruct nestStruct `json:"nest"`
		}

		full := request{
			Name:   "testerson",
			Age:    21,
			Active: false,
			NestedStruct: nestStruct{
				Test: "value",
			},
		}

		b, _ := json.Marshal(full)
		testStruct := request{}

		DecodeRequest(b, []string{"Name", "Age", "Active", "NestedStruct"}, &testStruct)
		if testStruct != full {
			t.Fail()
		}
	})
}
