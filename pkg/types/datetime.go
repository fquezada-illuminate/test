package types

import (
	"errors"
	"fmt"
	"gopkg.in/guregu/null.v3"
	"time"
)

type NullDatetime struct {
	null.Time
}
type NullDate struct {
	null.Time
}
type Datetime struct {
	time.Time
}
type Date struct {
	time.Time
}

const (
	FORMAT_DATETIME_INPUT  = "2006-01-02 15:04:05"
	FORMAT_DATETIME_OUTPUT = "2006-01-02T15:04:05-07:00"
	FORMAT_DATE            = "2006-01-02"
)

func parseTimeFromFormat(format string, input string, nullable bool) (*time.Time, error) {
	// Ignore null, like in the main JSON package.
	if input == "null" && nullable {
		return nil, nil
	}

	t, err := time.Parse(`"`+format+`"`, input)
	if err != nil {
		return nil, err
	}

	t.UTC()
	if t.Year() <= 0 || t.Year() > 10000 {
		return nil, errors.New("time year is not a valid year value, must be greater than 1 and less than 10000")
	}

	return &t, nil
}

func formatNullTime(nt null.Time, format string) ([]byte, error) {
	if nt.Valid == false {
		return []byte("null"), nil
	}

	return []byte(fmt.Sprintf(`"%s"`, nt.Time.Format(format))), nil
}

func formatTime(t time.Time, format string) ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.Format(format))), nil
}

// MarshalJSON will output the time in the form YYYY-MM-DD or null if applicable
func (nd NullDate) MarshalJSON() ([]byte, error) {
	return formatNullTime(nd.Time, FORMAT_DATE)
}

// UnmarshalJSON checks specifically for the time format YYYY-MM-DD or if its null return as such
func (nd *NullDate) UnmarshalJSON(data []byte) error {
	t, err := parseTimeFromFormat(FORMAT_DATE, string(data), true)
	if err != nil {
		return err
	}

	if t != nil {
		nd.SetValid(*t)
	} else {
		nd.Valid = false
		nd.Time.Time = time.Time{}
	}

	return nil
}

// MarshalJSON will output the time in the form YYYY-MM-DD HH:ii:ss or null if applicable
func (ndt NullDatetime) MarshalJSON() ([]byte, error) {
	return formatNullTime(ndt.Time, FORMAT_DATETIME_OUTPUT)
}

// UnmarshalJSON checks specifically for the time format YYYY-MM-DD HH:ii:ss or if its null return as such
func (ndt *NullDatetime) UnmarshalJSON(data []byte) error {
	t, err := parseTimeFromFormat(FORMAT_DATETIME_INPUT, string(data), true)
	if err != nil {
		return err
	}

	if t != nil {
		ndt.SetValid(*t)
	} else {
		ndt.Valid = false
		ndt.Time.Time = time.Time{}
	}

	return nil
}

//MarshalJSON will marshal the date into the format YYYY-MM-DD hh:ii:ss
func (dt Datetime) MarshalJSON() ([]byte, error) {
	return formatTime(dt.Time, FORMAT_DATETIME_OUTPUT)
}

//Unmarshal will attempt to parse the input into a YYYY-MM-DD hh:ii:ss format
func (dt *Datetime) UnmarshalJSON(data []byte) error {
	t, err := parseTimeFromFormat(FORMAT_DATETIME_INPUT, string(data), false)
	if err != nil {
		return err
	}

	dt.Time = *t
	return nil
}

//MarshalJSON will marshal the date into the format YYYY-MM-DD
func (d Date) MarshalJSON() ([]byte, error) {
	return formatTime(d.Time, FORMAT_DATE)
}

//Unmarshal will attempt to parse the input into a YYYY-MM-DD format
func (d *Date) UnmarshalJSON(data []byte) error {
	t, err := parseTimeFromFormat(FORMAT_DATE, string(data), false)
	if err != nil {
		return err
	}

	d.Time = *t
	return nil
}
