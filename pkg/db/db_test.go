package db

import "testing"

type temp struct {
	Field int `json:"fieldOne" db:"field_1"`
}

func TestGetJsonToDbMap(t *testing.T) {
	actual := GetJsonToDbMap(temp{})
	v, ok :=  actual["fieldOne"]
	if !ok {
		t.Errorf("The expected key '%s' does not exist.","fieldOne")
	}

	if v != "field_1"{
		t.Errorf("The expected value of was '%s': Got '%s'.","field_1", v)
	}
}

func TestGetColumns(t *testing.T) {
	actual := GetColumns(temp{})

	if len(actual) != 1 {
		t.Errorf("Expected to get an array of values that is length 1, but the length was  %d", len(actual))
	}

	if actual[0] != "field_1" {
		t.Errorf("Expected the first value to be 'field_1'. Got '%s'", actual[0])
	}
}

func TestStructToMap(t *testing.T) {
	type nestedTemp struct {
		Field string
		Nested struct {
			NestedField string
		}
	}

	actual := StructToMap(nestedTemp{})
	_, ok := actual["Field"]
	if !ok {
		t.Errorf("The expected key 'Field' does not exist.")
	}

	v, _ := actual["Nested"]
	nestedVal, ok := (v).(map[string]interface{})
	if !ok {
		t.Errorf("The expected key 'NestedField' does not isn't a map[string]interface.")

	}

	_, ok = nestedVal["NestedField"]
	if !ok {
		t.Errorf("The expected key nested value 'NestedField' does not exist.")
	}
}