package request

import (
	"fmt"
	"net/url"
	"strconv"
)

// Query contains the query string values and helper functions for type conversion
type Query struct {
	url.Values
}

func (q *Query) get(key string) ([]string, error) {

	values := q.Values[key]

	if len(values) < 1 {
		return values, fmt.Errorf("Missing key: %s", key)
	}

	return values, nil
}

// String returns the first value for the given key as a string
func (q *Query) String(key string) (string, error) {
	value, err := q.get(key)

	if err != nil {
		return "", err
	}

	return value[0], nil
}

// Strings returns all values for the given key as strings
func (q *Query) Strings(key string) ([]string, error) {
	values, err := q.get(key)

	if err != nil {
		return []string{}, err
	}

	return values, nil
}

// Int returns the first value for the given key as an int (int32)
func (q *Query) Int(key string) (int, error) {
	values, err := q.get(key)

	if err != nil {
		return 0, err
	}

	i, err := strconv.Atoi(values[0])

	if err != nil {
		return 0, err
	}

	return i, nil
}

// Ints returns all values for the given key as ints (int32)
func (q *Query) Ints(key string) ([]int, error) {
	values, err := q.get(key)

	if err != nil {
		return []int{}, err
	}

	ints := make([]int, len(values))

	for index, i := range values {
		ints[index], err = strconv.Atoi(i)

		if err != nil {
			return []int{}, err
		}
	}

	return ints, nil
}

// Float returns the first value for the given key as a float64
func (q *Query) Float(key string) (float64, error) {
	values, err := q.get(key)

	if err != nil {
		return 0, err
	}

	f, err := strconv.ParseFloat(values[0], 64)

	if err != nil {
		return 0, err
	}

	return f, nil
}

// Floats returns all values for the given key as float64
func (q *Query) Floats(key string) ([]float64, error) {
	values, err := q.get(key)

	if err != nil {
		return []float64{}, err
	}

	floats := make([]float64, len(values))

	for index, f := range values {
		floats[index], err = strconv.ParseFloat(f, 64)

		if err != nil {
			return []float64{}, err
		}
	}

	return floats, nil
}

func toBool(s string) (bool, error) {
	if s == "true" {
		return true, nil
	}

	if s == "false" {
		return false, nil
	}

	return false, fmt.Errorf("Can not convert %s to a boolean", s)
}

// Bool returns the first value for the given key as a boolean
func (q *Query) Bool(key string) (bool, error) {
	values, err := q.get(key)

	if err != nil {
		return false, err
	}

	b, err := toBool(values[0])

	if err != nil {
		return false, err
	}

	return b, nil
}

// Bools returns all values for the given key as booleans
func (q *Query) Bools(key string) ([]bool, error) {
	values, err := q.get(key)

	if err != nil {
		return []bool{}, err
	}

	bools := make([]bool, len(values))

	for i, b := range values {
		bools[i], err = toBool(b)

		if err != nil {
			return []bool{}, err
		}
	}

	return bools, nil
}
