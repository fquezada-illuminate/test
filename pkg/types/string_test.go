package types

import (
	"fmt"
	"testing"

	null "gopkg.in/guregu/null.v3"
)

// func NewString(s string, valid bool) NullString {
// 	return NullString{String: null.NewString(s, valid)}
// }

func TestNewString(t *testing.T) {
	tests := []struct {
		Name   string
		String string
		Valid  bool
	}{
		{"empty string; valid = true", "", true},
		{"empty string; valid = false", "", true},
		{"random string; valid = true", "random", true},
		{"random string; valid = false", "random", false},
	}

	for _, test := range tests {
		nns := null.NewString(test.String, test.Valid)
		tns := NewNullString(test.String, test.Valid)
		if nns != tns.String {
			t.Errorf("%s expected %v, got %v", test.Name, nns, tns)
		}
	}
}

func TestUnmarshalJSON(t *testing.T) {

	tests := []struct {
		Name           string
		Input          string
		ExpectedResult string
	}{
		// tests that should fail
		{"Unmarshal empty object", "{}", "json: cannot unmarshal {} into Go value of type types.NullString"},
		{"Unmarshal arbitrary object", `{"val": 1}`, `json: cannot unmarshal {"val": 1} into Go value of type types.NullString`},
		{"Unmarshal arbitraty object with null value", `{"val": null}`, `json: cannot unmarshal {"val": null} into Go value of type types.NullString`},
		{"Unmarshal empty array", `[]`, `json: cannot unmarshal [] into Go value of type types.NullString`},
		{"Unmarshal random array array", `[1, "something", true]`, `json: cannot unmarshal [1, "something", true] into Go value of type types.NullString`},
		{"Unmarshal bool", `true`, `json: cannot unmarshal bool into Go value of type null.String`},
		{"Unmarshal number", `17`, `json: cannot unmarshal float64 into Go value of type null.String`},
		{"Unmarshal error", ``, `unexpected end of JSON input`},
		// tests that should succeed
		{"Unmarshal null", "null", fmt.Sprintf("%v", nil)},
		{"Unmarshal empty string", `""`, ""},
	}

	for _, test := range tests {
		var ns NullString
		err := ns.UnmarshalJSON([]byte(test.Input))

		if err != nil && err.Error() != test.ExpectedResult {
			t.Errorf(`%s expected "%s", got "%v"`, test.Name, test.ExpectedResult, err)
		} else if val, _ := ns.Value(); err == nil && fmt.Sprintf("%v", val) != test.ExpectedResult {
			t.Errorf(`%s expected "%s", got "%v"`, test.Name, test.ExpectedResult, val)
		}
	}

}
