package request

import (
	"fmt"
	"net/http/httptest"
	"strconv"
	"testing"
)

const testKey = "test"

func TestGet(t *testing.T) {
	testCases := []struct {
		desc        string
		queryString string
		expected    []string
		err         error
	}{
		{
			desc:        "missing key",
			queryString: "otherkey=something",
			expected:    []string{},
			err:         fmt.Errorf("Missing key: %s", testKey),
		}, {
			desc:        "single value",
			queryString: fmt.Sprintf("%s=%s", testKey, "testval"),
			expected:    []string{"testval"},
			err:         nil,
		}, {
			desc:        "multiple of same key with different values",
			queryString: fmt.Sprintf("%s=%s&%s=%s&%s=%s", testKey, "testval", testKey, "testval1", testKey, "testval2"),
			expected:    []string{"testval", "testval1", "testval2"},
			err:         nil,
		}, {
			desc:        "multiple of same key with same value",
			queryString: fmt.Sprintf("%s=%s&%s=%s&%s=%s", testKey, "testval", testKey, "testval", testKey, "testval"),
			expected:    []string{"testval", "testval", "testval"},
			err:         nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			request := httptest.NewRequest("GET", fmt.Sprintf("/test?%s", tC.queryString), nil)
			qs := Query{Values: request.URL.Query()}

			values, err := qs.get(testKey)

			for i := range tC.expected {
				if tC.expected[i] != values[i] {
					t.Errorf("expected: %v, got %v", tC.expected[i], values[i])
				}
			}

			if tC.err != nil && err == nil {
				t.Errorf("expected error %v, got no error", tC.err)
			} else if tC.err == nil && err != nil {
				t.Errorf("did not expect an error, got %v", err)
			} else if tC.err != nil && err != nil && err.Error() != tC.err.Error() {
				t.Errorf("expected error %v, got %v", tC.err, err)
			}
		})
	}
}

func TestString(t *testing.T) {
	testCases := []struct {
		desc        string
		queryString string
		expected    string
		err         error
	}{
		{
			desc:        "single value",
			queryString: fmt.Sprintf("%s=%s", testKey, "test"),
			expected:    "test",
			err:         nil,
		}, {
			desc:        "multiple values",
			queryString: fmt.Sprintf("%s=%s&%s=%s&%s=%s&", testKey, "test3", testKey, "test2", testKey, "test1"),
			expected:    "test3",
			err:         nil,
		}, {
			desc:        "encoded value",
			queryString: fmt.Sprintf("%s=%s", testKey, "test%20text%3F"),
			expected:    "test text?",
			err:         nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			request := httptest.NewRequest("GET", fmt.Sprintf("/test?%s", tC.queryString), nil)
			qs := Query{Values: request.URL.Query()}

			value, err := qs.String(testKey)

			if value != tC.expected {
				t.Errorf("expected %v got %v", tC.expected, value)
			}

			if tC.err != nil && err == nil {
				t.Errorf("expected error %v, got no error", tC.err)
			} else if tC.err == nil && err != nil {
				t.Errorf("did not expect an error, got %v", err)
			} else if tC.err != nil && err != nil && err.Error() != tC.err.Error() {
				t.Errorf("expected error %v, got %v", tC.err, err)
			}
		})
	}
}

func TestStrings(t *testing.T) {
	testCases := []struct {
		desc        string
		queryString string
		expected    []string
		err         error
	}{
		{
			desc:        "single value",
			queryString: fmt.Sprintf("%s=%s", testKey, "test"),
			expected:    []string{"test"},
			err:         nil,
		}, {
			desc:        "multiple values",
			queryString: fmt.Sprintf("%s=%s&%s=%s&%s=%s&", testKey, "test1", testKey, "test2", testKey, "test3"),
			expected:    []string{"test1", "test2", "test3"},
			err:         nil,
		}, {
			desc:        "encoded value",
			queryString: fmt.Sprintf("%s=%s", testKey, "test%20text%3F"),
			expected:    []string{"test text?"},
			err:         nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			request := httptest.NewRequest("GET", fmt.Sprintf("/test?%s", tC.queryString), nil)
			qs := Query{Values: request.URL.Query()}

			values, err := qs.Strings(testKey)

			for i := 0; i < len(tC.expected) && i < len(values); i++ {
				if tC.expected[i] != values[i] {
					t.Errorf("expected: %v, got %v", tC.expected[i], values[i])
				}
			}

			if tC.err != nil && err == nil {
				t.Errorf("expected error %v, got no error", tC.err)
			} else if tC.err == nil && err != nil {
				t.Errorf("did not expect an error, got %v", err)
			} else if tC.err != nil && err != nil && err.Error() != tC.err.Error() {
				t.Errorf("expected error %v, got %v", tC.err, err)
			}
		})
	}
}

func TestInt(t *testing.T) {
	testCases := []struct {
		desc        string
		queryString string
		expected    int
		err         error
	}{
		{
			desc:        "single value",
			queryString: fmt.Sprintf("%s=%d", testKey, 12),
			expected:    12,
			err:         nil,
		}, {
			desc:        "multiple values",
			queryString: fmt.Sprintf("%s=%d&%s=%d&%s=%d", testKey, 42, testKey, 12, testKey, 13),
			expected:    42,
			err:         nil,
		}, {
			desc:        "invalid single value",
			queryString: fmt.Sprintf("%s=not_an_int", testKey),
			expected:    0,
			err:         fmt.Errorf((&strconv.NumError{Func: "Atoi", Num: "not_an_int", Err: strconv.ErrSyntax}).Error()),
		}, {
			desc:        "invalid multiple values (first value)",
			queryString: fmt.Sprintf("%s=%s&%s=%d&%s=%d", testKey, "not_an_int", testKey, 42, testKey, 13),
			expected:    0,
			err:         fmt.Errorf((&strconv.NumError{Func: "Atoi", Num: "not_an_int", Err: strconv.ErrSyntax}).Error()),
		}, {
			desc:        "invalid multiple values (not first value)",
			queryString: fmt.Sprintf("%s=%d&%s=%s&%s=%d", testKey, 42, testKey, "not_an_int", testKey, 13),
			expected:    42,
			err:         nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			request := httptest.NewRequest("GET", fmt.Sprintf("/test?%s", tC.queryString), nil)
			qs := Query{Values: request.URL.Query()}

			value, err := qs.Int(testKey)

			if tC.expected != value {
				t.Errorf("expected: %v, got %v", tC.expected, value)
			}

			if tC.err != nil && err == nil {
				t.Errorf("expected error %v, got no error", tC.err)
			} else if tC.err == nil && err != nil {
				t.Errorf("did not expect an error, got %v", err)
			} else if tC.err != nil && err != nil && err.Error() != tC.err.Error() {
				t.Errorf("expected error %v, got %v", tC.err, err)
			}
		})
	}
}

func TestInts(t *testing.T) {
	testCases := []struct {
		desc        string
		queryString string
		expected    []int
		err         error
	}{
		{
			desc:        "single value",
			queryString: fmt.Sprintf("%s=%d", testKey, 12),
			expected:    []int{12},
			err:         nil,
		}, {
			desc:        "multiple values",
			queryString: fmt.Sprintf("%s=%d&%s=%d&%s=%d", testKey, 42, testKey, 12, testKey, 13),
			expected:    []int{42, 12, 13},
			err:         nil,
		}, {
			desc:        "invalid single value",
			queryString: fmt.Sprintf("%s=not_an_int", testKey),
			expected:    []int{12},
			err:         fmt.Errorf((&strconv.NumError{Func: "Atoi", Num: "not_an_int", Err: strconv.ErrSyntax}).Error()),
		}, {
			desc:        "invalid multiple values",
			queryString: fmt.Sprintf("%s=%d&%s=%s&%s=%d", testKey, 42, testKey, "not_an_int", testKey, 13),
			expected:    []int{12},
			err:         fmt.Errorf((&strconv.NumError{Func: "Atoi", Num: "not_an_int", Err: strconv.ErrSyntax}).Error()),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			request := httptest.NewRequest("GET", fmt.Sprintf("/test?%s", tC.queryString), nil)
			qs := Query{Values: request.URL.Query()}

			values, err := qs.Ints(testKey)

			for i := 0; i < len(tC.expected) && i < len(values); i++ {
				if tC.expected[i] != values[i] {
					t.Errorf("expected: %v, got %v", tC.expected[i], values[i])
				}
			}

			if tC.err != nil && err == nil {
				t.Errorf("expected error %v, got no error", tC.err)
			} else if tC.err == nil && err != nil {
				t.Errorf("did not expect an error, got %v", err)
			} else if tC.err != nil && err != nil && err.Error() != tC.err.Error() {
				t.Errorf("expected error %v, got %v", tC.err, err)
			}
		})
	}
}

func TestFloat(t *testing.T) {
	testCases := []struct {
		desc        string
		queryString string
		expected    float64
		err         error
	}{
		{
			desc:        "single value",
			queryString: fmt.Sprintf("%s=%s", testKey, "12.1"),
			expected:    12.1,
			err:         nil,
		}, {
			desc:        "multiple values",
			queryString: fmt.Sprintf("%s=%s&%s=%s&%s=%s", testKey, "42.1234", testKey, "12.1", testKey, "13.37"),
			expected:    42.1234,
			err:         nil,
		}, {
			desc:        "invalid single value",
			queryString: fmt.Sprintf("%s=not_a_float", testKey),
			expected:    0,
			err:         fmt.Errorf((&strconv.NumError{Func: "ParseFloat", Num: "not_a_float", Err: strconv.ErrSyntax}).Error()),
		}, {
			desc:        "invalid multiple values (first value)",
			queryString: fmt.Sprintf("%s=%s&%s=%s&%s=%s", testKey, "not_a_float", testKey, "42.1234", testKey, "13.37"),
			expected:    0,
			err:         fmt.Errorf((&strconv.NumError{Func: "ParseFloat", Num: "not_a_float", Err: strconv.ErrSyntax}).Error()),
		}, {
			desc:        "invalid multiple values (not first value)",
			queryString: fmt.Sprintf("%s=%s&%s=%s&%s=%s", testKey, "42.1234", testKey, "not_a_float", testKey, "13.37"),
			expected:    42.1234,
			err:         nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			request := httptest.NewRequest("GET", fmt.Sprintf("/test?%s", tC.queryString), nil)
			qs := Query{Values: request.URL.Query()}

			value, err := qs.Float(testKey)

			if tC.expected != value {
				t.Errorf("expected: %v, got %v", tC.expected, value)
			}

			if tC.err != nil && err == nil {
				t.Errorf("expected error %v, got no error", tC.err)
			} else if tC.err == nil && err != nil {
				t.Errorf("did not expect an error, got %v", err)
			} else if tC.err != nil && err != nil && err.Error() != tC.err.Error() {
				t.Errorf("expected error %v, got %v", tC.err, err)
			}
		})
	}
}

func TestFloats(t *testing.T) {
	testCases := []struct {
		desc        string
		queryString string
		expected    []float64
		err         error
	}{
		{
			desc:        "single value",
			queryString: fmt.Sprintf("%s=%s", testKey, "12.1"),
			expected:    []float64{12.1},
			err:         nil,
		}, {
			desc:        "multiple values",
			queryString: fmt.Sprintf("%s=%s&%s=%s&%s=%s", testKey, "42.1234", testKey, "12.1", testKey, "13.37"),
			expected:    []float64{42.1234, 12.1, 13.37},
			err:         nil,
		}, {
			desc:        "invalid single value",
			queryString: fmt.Sprintf("%s=not_a_float", testKey),
			expected:    []float64{},
			err:         fmt.Errorf((&strconv.NumError{Func: "ParseFloat", Num: "not_a_float", Err: strconv.ErrSyntax}).Error()),
		}, {
			desc:        "invalid multiple values",
			queryString: fmt.Sprintf("%s=%s&%s=%s&%s=%s", testKey, "42.1234", testKey, "not_a_float", testKey, "13.37"),
			expected:    []float64{},
			err:         fmt.Errorf((&strconv.NumError{Func: "ParseFloat", Num: "not_a_float", Err: strconv.ErrSyntax}).Error()),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			request := httptest.NewRequest("GET", fmt.Sprintf("/test?%s", tC.queryString), nil)
			qs := Query{Values: request.URL.Query()}

			values, err := qs.Floats(testKey)

			for i := 0; i < len(tC.expected) && i < len(values); i++ {
				if tC.expected[i] != values[i] {
					t.Errorf("expected: %v, got %v", tC.expected[i], values[i])
				}
			}

			if tC.err != nil && err == nil {
				t.Errorf("expected error %v, got no error", tC.err)
			} else if tC.err == nil && err != nil {
				t.Errorf("did not expect an error, got %v", err)
			} else if tC.err != nil && err != nil && err.Error() != tC.err.Error() {
				t.Errorf("expected error %v, got %v", tC.err, err)
			}
		})
	}
}

func TestToBool(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected bool
		err      error
	}{
		{
			desc:     "true",
			input:    "true",
			expected: true,
			err:      nil,
		}, {
			desc:     "false",
			input:    "false",
			expected: false,
			err:      nil,
		}, {
			desc:     "True",
			input:    "True",
			expected: false,
			err:      fmt.Errorf("Can not convert %s to a boolean", "True"),
		}, {
			desc:     "False",
			input:    "False",
			expected: false,
			err:      fmt.Errorf("Can not convert %s to a boolean", "False"),
		}, {
			desc:     "invalid value",
			input:    "not_a_bool",
			expected: false,
			err:      fmt.Errorf("Can not convert %s to a boolean", "not_a_bool"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			value, err := toBool(tC.input)

			if tC.expected != value {
				t.Errorf("expected: %v, got %v", tC.expected, value)
			}

			if tC.err != nil && err == nil {
				t.Errorf("expected error %v, got no error", tC.err)
			} else if tC.err == nil && err != nil {
				t.Errorf("did not expect an error, got %v", err)
			} else if tC.err != nil && err != nil && err.Error() != tC.err.Error() {
				t.Errorf("expected error %v, got %v", tC.err, err)
			}
		})
	}
}

func TestBool(t *testing.T) {
	testCases := []struct {
		desc        string
		queryString string
		expected    bool
		err         error
	}{
		{
			desc:        "single value",
			queryString: fmt.Sprintf("%s=%s", testKey, "true"),
			expected:    true,
			err:         nil,
		}, {
			desc:        "multiple values",
			queryString: fmt.Sprintf("%s=%s&%s=%s&%s=%s&%s=%s", testKey, "false", testKey, "true", testKey, "true", testKey, "false"),
			expected:    false,
			err:         nil,
		}, {
			desc:        "single invalid value",
			queryString: fmt.Sprintf("%s=%s", testKey, "not_a_bool"),
			expected:    false,
			err:         fmt.Errorf("Can not convert %s to a boolean", "not_a_bool"),
		}, {
			desc:        "invalid multiple values (not first one)",
			queryString: fmt.Sprintf("%s=%s&%s=%s", testKey, "true", testKey, "not_a_bool"),
			expected:    true,
			err:         nil,
		}, {
			desc:        "invalid multiple values (first one)",
			queryString: fmt.Sprintf("%s=%s&%s=%s", testKey, "not_a_bool", testKey, "true"),
			expected:    false,
			err:         fmt.Errorf("Can not convert %s to a boolean", "not_a_bool"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			request := httptest.NewRequest("GET", fmt.Sprintf("/test?%s", tC.queryString), nil)
			qs := Query{Values: request.URL.Query()}

			value, err := qs.Bool(testKey)

			if tC.expected != value {
				t.Errorf("expected: %v, got %v", tC.expected, value)
			}

			if tC.err != nil && err == nil {
				t.Errorf("expected error %v, got no error", tC.err)
			} else if tC.err == nil && err != nil {
				t.Errorf("did not expect an error, got %v", err)
			} else if tC.err != nil && err != nil && err.Error() != tC.err.Error() {
				t.Errorf("expected error %v, got %v", tC.err, err)
			}
		})
	}
}

func TestBools(t *testing.T) {
	testCases := []struct {
		desc        string
		queryString string
		expected    []bool
		err         error
	}{
		{
			desc:        "single value",
			queryString: fmt.Sprintf("%s=%s", testKey, "true"),
			expected:    []bool{true},
			err:         nil,
		}, {
			desc:        "multiple values",
			queryString: fmt.Sprintf("%s=%s&%s=%s&%s=%s", testKey, "true", testKey, "false", testKey, "false"),
			expected:    []bool{true, false, false},
			err:         nil,
		}, {
			desc:        "invalid single value",
			queryString: fmt.Sprintf("%s=not_a_bool", testKey),
			expected:    []bool{},
			err:         fmt.Errorf("Can not convert %s to a boolean", "not_a_bool"),
		}, {
			desc:        "invalid multiple values",
			queryString: fmt.Sprintf("%s=%s&%s=%s&%s=%s", testKey, "false", testKey, "not_a_bool", testKey, "true"),
			expected:    []bool{},
			err:         fmt.Errorf("Can not convert %s to a boolean", "not_a_bool"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			request := httptest.NewRequest("GET", fmt.Sprintf("/test?%s", tC.queryString), nil)
			qs := Query{Values: request.URL.Query()}

			values, err := qs.Bools(testKey)

			for i := 0; i < len(tC.expected) && i < len(values); i++ {
				if tC.expected[i] != values[i] {
					t.Errorf("expected: %v, got %v", tC.expected[i], values[i])
				}
			}

			if tC.err != nil && err == nil {
				t.Errorf("expected error %v, got no error", tC.err)
			} else if tC.err == nil && err != nil {
				t.Errorf("did not expect an error, got %v", err)
			} else if tC.err != nil && err != nil && err.Error() != tC.err.Error() {
				t.Errorf("expected error %v, got %v", tC.err, err)
			}
		})
	}
}
