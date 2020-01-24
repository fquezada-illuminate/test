package structs

import (
	"fmt"
	"reflect"
	"testing"
)

type MockStruct struct {
	Id          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	DisplayName string `json:"display_name" db:"display_name"`
	Mixed       string `db:"mixed"`
}

func TestGetTagMapErrors(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			key   string
			value string
		}
		expectedResponse map[string]string
	}{
		{
			name: "tag for key does not exist",
			input: struct {
				key   string
				value string
			}{key: "invalid", value: "json"},
			expectedResponse: map[string]string{},
		},
		{
			name: "tag for value does not exist",
			input: struct {
				key string
				value string
			}{key: "db", value: "invalid"},
			expectedResponse: map[string]string{},
		},
		{
			name: "tag for key does not exist on 1 field",
			input: struct {
				key   string
				value string
			}{key: "json", value: "db"},
			expectedResponse: map[string]string {
				"id": "id",
				"name": "name",
				"display_name": "display_name",
			},
		},
	}

	sh := Helper{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sMap, err := sh.GetTagMap(MockStruct{}, test.input.key, test.input.value)
			if err != nil {
				t.Errorf("Expected response and got error: %s", test.name)
			}

			if !reflect.DeepEqual(test.expectedResponse, sMap) {
				fmt.Println(test.expectedResponse, sMap)
				t.Error("Unexpected map result")
			}
		})
	}
}

func TestGetTagMapNoError(t *testing.T) {

	goodMockStruct := struct {
		Id          string `json:"id" db:"id"`
		Name        string `json:"name" db:"name"`
		DisplayName string `json:"display_name" db:"display_name_test"`
	}{}

	sh := Helper{}

	sMap, err := sh.GetTagMap(goodMockStruct, "json", "db")
	if err != nil {
		t.Errorf("Expected response and got error %s", err)
	}

	expectedMap := map[string]string{
		"id":           "id",
		"name":         "name",
		"display_name": "display_name_test",
	}

	if !reflect.DeepEqual(expectedMap, sMap) {
		t.Error("Unexpected map result")
	}
}

func TestGetTagValues(t *testing.T) {

	sh := Helper{}
	expectedMap := []string{
		"id",
		"name",
		"display_name",
		"mixed",
	}

	values := sh.GetTagValues(MockStruct{}, "db")

	if !reflect.DeepEqual(expectedMap, values) {
		t.Error("Unexpected slice result")
	}
}

func TestGetTagValuesNotAllFieldsHaveTag(t *testing.T) {
	sh := Helper{}

	testMockStruct := struct {
		Id          string `json:"id" db:"id"`
		Name        string `json:"name" db:"name"`
		DisplayName string `json:"display_name"`
		Mixed       string `json:"mixed" db:"-"`
	}{}

	expectedMap := []string{
		"id",
		"name",
	}


	values := sh.GetTagValues(testMockStruct, "db")

	if !reflect.DeepEqual(expectedMap, values) {
		t.Error("Unexpected slice result")
	}
}

func TestGetMapByTag(t *testing.T) {
	sh := Helper{}
	mock := MockStruct{
		"1",
		"Mock",
		"Mock Struct",
		"asdf",
	}
	expectedMap := map[string]interface{}{
		"id":           mock.Id,
		"name":         mock.Name,
		"display_name": mock.DisplayName,
		"mixed":        mock.Mixed,
	}

	sMap := sh.GetMapByTag(mock, "db")

	if !reflect.DeepEqual(expectedMap, sMap) {
		t.Error("Unexpected map result")
	}
}
