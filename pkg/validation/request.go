package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// DecodeRequest will take a json request body and validate that all the top-level fields are in members of the struct
// and that each key is within the validaFields array.  It will match the request parameter by its json tag name if it
// exists and it will match the valid fields By the name of the struct member.
func DecodeRequest(body []byte, validFields []string, objPtr interface{}) error {
	rv := reflect.ValueOf(objPtr)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		panic(errors.New("obj should be a pointer to a struct"))
	}

	objType := rv.Elem().Type()
	for _, fieldName := range validFields {
		_, ok := objType.FieldByName(fieldName)
		if !ok {
			panic(fmt.Sprintf("%s is not a part of %s", fieldName, objType.String()))
		}
	}

	requestValues := make(map[string]json.RawMessage)
	err := json.Unmarshal(body, &requestValues)
	if err != nil {
		return err
	}

	fieldMap := getFieldsMappedToJsonTag(objType)

	errs := make([]string, 0)
	for jsonKey := range requestValues {
		if _, ok := fieldMap[jsonKey]; !ok {
			errs = append(errs, fmt.Sprintf("%s: This property does not exist.", jsonKey))
			continue
		}

		// don't add an additional validation checks if the array is empty
		if len(validFields) == 0 {
			continue
		}

		field := fieldMap[jsonKey]
		canBeSet := false
		for _, f := range validFields {
			if f == field {
				canBeSet = true
				break
			}
		}

		if !canBeSet {
			errs = append(errs, fmt.Sprintf("%s: This property is not allowed to be set.", jsonKey))
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ERROR_DELIMITER))
	}

	// decode each key in the request individually
	for jsonKey, rawJson := range requestValues {
		field := rv.Elem().FieldByName(fieldMap[jsonKey])
		fieldValuePtr := reflect.New(field.Type()).Interface()

		err := json.Unmarshal(rawJson, fieldValuePtr)
		if err != nil {
			errs = append(errs, err.Error())
		} else {
			field.Set(reflect.ValueOf(fieldValuePtr).Elem())
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ERROR_DELIMITER))
	}

	return nil
}

// getFieldsMappedToJsonTag returns a map of fields where the key is the json tag or the field name if a json
// field name isn't set and the value is the real field name
func getFieldsMappedToJsonTag(structType reflect.Type) map[string]string {
	fields := make(map[string]string, structType.NumField())

	for i := 0; i < structType.NumField(); i++ {
		f := structType.Field(i)
		key := f.Name

		jsonTag, ok := f.Tag.Lookup("json")
		if ok {
			key = strings.Split(jsonTag, ",")[0]
		}

		fields[key] = f.Name
	}

	return fields
}
