package validation

import (
	"reflect"
	"testing"
)

func TestStringNotBlank(t *testing.T) {
	t.Run("String of some length should pass", func(t *testing.T) {
		str := "string"
		if !StringNotBlank(reflect.ValueOf(str), nil, nil) {
			t.Errorf("A non NullString value should return true but got false")
		}
	})

	t.Run("A nil value should pass", func(t *testing.T) {
		var v *string = nil
		if !StringNotBlank(reflect.ValueOf(v), nil, nil) {
			t.Errorf("A nil value should return true but got false")
		}
	})

	t.Run("An empty string returns false", func(t *testing.T) {
		str := ""
		if StringNotBlank(reflect.ValueOf(str), nil, nil) {
			t.Errorf("A populated NullString value should return true but got false")
		}
	})

	t.Run("A 0 returns true", func(t *testing.T) {
		intVal := 0
		if !StringNotBlank(reflect.ValueOf(intVal), nil, nil) {
			t.Errorf("A 0 value should return true, false returned")
		}
	})

	t.Run("A non-zero integer returns false", func(t *testing.T) {
		intVal := 2
		if StringNotBlank(reflect.ValueOf(intVal), nil, nil) {
			t.Errorf("A non-zero integer returns false, true returned")
		}
	})

	t.Run("A string that is shorter than min, returns false", func(t *testing.T) {
		str := "small"
		min := 6
		if StringNotBlank(reflect.ValueOf(str), &min, nil) {
			t.Errorf("A string shorter than the min returned true")
		}
	})

	t.Run("A string that is longer than max, returns false", func(t *testing.T) {
		str := "small"
		max := 2
		if StringNotBlank(reflect.ValueOf(str), nil, &max) {
			t.Errorf("A string longer than the min returned true")
		}
	})
}

func TestUuidNotBlank(t *testing.T) {
	t.Run("A null value should pass", func(t *testing.T) {
		var v *string = nil
		if !UuidNotBlank(reflect.ValueOf(v)) {
			t.Errorf("A nil value should return true but got false")
		}
	})

	t.Run("A valid UUIDv4 should pass", func(t *testing.T) {
		str := "f99798eb-6724-4448-9ee0-9259cde51c2a"
		if !UuidNotBlank(reflect.ValueOf(str)) {
			t.Errorf("A valid UUIDv4 should return true but got false")
		}
	})

	t.Run("A 0 returns true", func(t *testing.T) {
		intVal := 0
		if !UuidNotBlank(reflect.ValueOf(intVal)) {
			t.Errorf("A 0 value should return true but got false")
		}
	})

	t.Run("An empty string returns false", func(t *testing.T) {
		str := ""
		if UuidNotBlank(reflect.ValueOf(str)) {
			t.Errorf("An empty string should return false but got true")
		}
	})
}
