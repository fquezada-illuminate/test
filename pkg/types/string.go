package types

import (
	"encoding/json"
	"fmt"

	"gopkg.in/guregu/null.v3"
)

// NullString is a wrapper around gopkg.in/guregu/null.v3 null.String
type NullString struct {
	null.String
}

// NewString creates a new null string
func NewNullString(s string, valid bool) NullString {
	return NullString{String: null.NewString(s, valid)}
}

// UnmarshalJSON unmarshals data into a NullString if and only if data is a string
// it uses the underlying null.String for unmarshalling unless the JSON is malformed
// or the type is a JSON object
func (s *NullString) UnmarshalJSON(data []byte) error {

	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v.(type) {
	case map[string]interface{}:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type types.NullString", string(data))
	case []interface{}:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type types.NullString", string(data))
	default:
		return s.String.UnmarshalJSON(data)
	}
	return err
}
