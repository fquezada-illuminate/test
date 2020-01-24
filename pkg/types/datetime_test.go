package types

import (
	"fmt"
	"testing"
	"time"
)

func TestNullDatetime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name           string
		input          []byte
		expectingError bool
	}{
		{
			name:           "Poorly formatted string",
			input:          []byte(`"Monday 2006 11:15pm"`),
			expectingError: true,
		},
		{
			name:           "Poorly formatted string",
			input:          []byte(`"Not even a time"`),
			expectingError: true,
		},
		{
			name:           "Edge case datetime",
			input:          []byte(`"0000-01-20 15:01:00"`),
			expectingError: true,
		},
		{
			name:           "A valid datetime string",
			input:          []byte(`"2018-01-20 15:01:00"`),
			expectingError: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ndt := NullDatetime{}
			err := ndt.UnmarshalJSON(test.input)

			if (err != nil) != test.expectingError {
				t.Errorf("Did not get an error when one was expected: %s", err.Error())
			}

			if test.expectingError == false {
				if ndt.Valid != true {
					t.Error("A valid datetime was expected to set the Valid to true")
				}

				tm, _ := ndt.Value()
				if tm == (time.Time{}) {
					t.Error("A valid datetime did not set the time to a non-empty value")
				}

			}
		})
	}

	t.Run("A json null is a value input", func(t *testing.T) {
		ndt := NullDatetime{}
		err := ndt.UnmarshalJSON([]byte(`null`))

		if err != nil {
			t.Errorf("json nulls are valid, error returned: %s", err.Error())
		}

		if ndt.Valid != false {
			t.Error("If the null valid should be false")
		}

		tm := ndt.Time.Time
		if tm != (time.Time{}) {
			t.Error("If the null, the time value should be empty")
		}
	})
}

func TestNullDatetime_MarshalJSON(t *testing.T) {
	t.Run("A null value should return as a 'null' string", func(t *testing.T) {
		ndt := NullDatetime{}

		b, err := ndt.MarshalJSON()
		if err != nil {
			t.Error("A null NullDateTime should be able to be marshaled without error")
		}

		if string(b) != "null" {
			t.Errorf("Expected marshaled null is incorrect: Got %s expected 'null", string(b))
		}
	})

	t.Run("Check a valid datetime", func(t *testing.T) {
		expectedTimeString := "2018-01-02 03:40:56"
		expectedTime, _ := time.Parse(FORMAT_DATETIME_INPUT, expectedTimeString)

		ndt := NullDatetime{}
		ndt.SetValid(expectedTime)

		b, err := ndt.MarshalJSON()
		if err != nil {
			t.Error("A valid time should be able to be marshaled without error")
		}

		// wrap expected results in quotes
		expectedTimeString = fmt.Sprintf(`"%s"`, expectedTime.Format(FORMAT_DATETIME_OUTPUT))
		if string(b) != expectedTimeString {
			t.Errorf("Expected marshaled string is incorrect: Got %s expected %s", string(b), expectedTimeString)
		}
	})
}

func TestNullDate_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name           string
		input          []byte
		expectingError bool
	}{
		{
			name:           "Poorly formatted string",
			input:          []byte(`"Monday 2006 11:15pm"`),
			expectingError: true,
		},
		{
			name:           "A valid time stamp with more granularity than expected",
			input:          []byte(`"2018-01-20 15:01:00"`),
			expectingError: true,
		},
		{
			name:           "Edge case datetime",
			input:          []byte(`"0000-01-20"`),
			expectingError: true,
		},
		{
			name:           "A valid datetime string",
			input:          []byte(`"2018-01-20"`),
			expectingError: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nd := NullDate{}
			err := nd.UnmarshalJSON(test.input)

			if (err != nil) != test.expectingError {
				t.Error("Did not get an error when one was expected")
			}

			if test.expectingError == false {
				if nd.Valid != true {
					t.Error("A valid datetime was expected to set the Valid to true")
				}

				tm, _ := nd.Value()
				if tm == (time.Time{}) {
					t.Error("A valid datetime did not set the time to a non-empty value")
				}

			}
		})
	}

	t.Run("A json null is a value input", func(t *testing.T) {
		nd := NullDate{}
		err := nd.UnmarshalJSON([]byte(`null`))

		if err != nil {
			t.Errorf("json nulls are valid, error returned: %s", err.Error())
		}

		if nd.Valid != false {
			t.Error("If the null valid should be false")
		}

		tm, _ := nd.Value()
		if tm == (time.Time{}) {
			t.Error("If the null, the time value should be empty")
		}
	})
}

func TestNullDate_MarshalJSON(t *testing.T) {
	t.Run("A null value should return as a 'null' string", func(t *testing.T) {
		nd := NullDate{}

		b, err := nd.MarshalJSON()
		if err != nil {
			t.Error("A null NullDateTime should be able to be marshaled without error")
		}

		if string(b) != "null" {
			t.Errorf("Expected marshaled null is incorrect: Got %s expected 'null'", string(b))
		}
	})

	t.Run("", func(t *testing.T) {
		expectedTimeString := "2018-01-02"
		expectedTime, _ := time.Parse(FORMAT_DATE, expectedTimeString)

		nd := NullDate{}
		nd.SetValid(expectedTime)
		b, err := nd.MarshalJSON()
		if err != nil {
			t.Error("A valid time should be able to be marshaled without error")
		}

		// wrap expected results in quotes
		expectedTimeString = fmt.Sprintf(`"%s"`, expectedTimeString)
		if string(b) != expectedTimeString {
			t.Errorf("Expected marshaled string is incorrect: Got %s expected %s", string(b), expectedTimeString)
		}
	})
}

func TestDatetime_MarshalJSON(t *testing.T) {
	t.Run("Marshalling a datetime will include provide the correct output", func(t *testing.T) {
		expectedDatetime := "2018-11-15 05:10:55"
		tm, _ := time.Parse(FORMAT_DATETIME_INPUT, expectedDatetime)
		dt := Datetime{Time: tm}

		actual, _ := dt.MarshalJSON()

		expectedDatetime = fmt.Sprintf(`"%s"`, tm.Format(FORMAT_DATETIME_OUTPUT))
		if string(actual) != expectedDatetime {
			t.Errorf("Expected Datetime unmarshal format is incorrect: Got %s expected %s", string(actual), expectedDatetime)
		}
	})
}

func TestDateTime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name           string
		input          []byte
		expectingError bool
	}{
		{
			name:           "Poorly formatted string",
			input:          []byte(`"Monday 2006 11:15pm"`),
			expectingError: true,
		},
		{
			name:           "Edge case datetime",
			input:          []byte(`"0000-01-20 13:14:50"`),
			expectingError: true,
		},
		{
			name:           "A valid datetime string",
			input:          []byte(`"2018-01-20 13:14:50"`),
			expectingError: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dt := Datetime{}
			err := dt.UnmarshalJSON(test.input)

			if (err != nil) != test.expectingError {
				t.Error("Did not get an error when one was expected")
			}

			if test.expectingError == false {
				if dt == (Datetime{}) {
					t.Error("A valid datetime did not set the time to a non-empty value")
				}

			}
		})
	}

	t.Run("A json null is a value input", func(t *testing.T) {
		nd := NullDate{}
		err := nd.UnmarshalJSON([]byte(`null`))

		if err != nil {
			t.Errorf("json nulls are valid, error returned: %s", err.Error())
		}

		if nd.Valid != false {
			t.Error("If the null valid should be false")
		}
		tm := nd.Time.Time
		if tm != (time.Time{}) {
			t.Error("If the null, the time value should be empty")
		}
	})
}

func TestDate_MarshalJSON(t *testing.T) {
	t.Run("Marshalling a Date will include provide the correct output", func(t *testing.T) {
		expectedDate := "2018-11-15"
		tm, _ := time.Parse(FORMAT_DATE, expectedDate)
		dt := Date{Time: tm}

		actual, _ := dt.MarshalJSON()

		expectedDate = fmt.Sprintf(`"%s"`, expectedDate)
		if string(actual) != expectedDate {
			t.Errorf("Expected Date unmarshal format is incorrect: Got %s expected %s", string(actual), expectedDate)
		}
	})
}

func TestDate_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name           string
		input          []byte
		expectingError bool
	}{
		{
			name:           "Invalid date format",
			input:          []byte(`"Monday 2006"`),
			expectingError: true,
		},
		{
			name:           "Edge case Date",
			input:          []byte(`"0000-01-20"`),
			expectingError: true,
		},
		{
			name:           "A valid Date string",
			input:          []byte(`"2018-01-20"`),
			expectingError: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := Date{}
			err := d.UnmarshalJSON(test.input)

			if (err != nil) != test.expectingError {
				t.Error("Did not get an error when one was expected")
			}

			if test.expectingError == false {
				if d == (Date{}) {
					t.Error("A valid Date did not set the time to a non-empty value")
				}

			}
		})
	}

	t.Run("A json null is a value input", func(t *testing.T) {
		nd := NullDate{}
		err := nd.UnmarshalJSON([]byte(`null`))

		if err != nil {
			t.Errorf("json nulls are valid, error returned: %s", err.Error())
		}

		if nd.Valid != false {
			t.Error("If the null valid should be false")
		}
		tm := nd.Time.Time
		if tm != (time.Time{}) {
			t.Error("If the null, the time value should be empty")
		}
	})
}
