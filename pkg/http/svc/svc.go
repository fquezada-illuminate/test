package svc

import (
	"encoding/json"
	"errors"
	"github.com/fatih/structs"
	"github.com/gorilla/mux"
	"github.com/illuminateeducation/rest-service-lib-go/pkg/db"
	"github.com/illuminateeducation/rest-service-lib-go/pkg/http/route"
	"github.com/illuminateeducation/rest-service-lib-go/pkg/response"
	"github.com/illuminateeducation/rest-service-lib-go/pkg/types"
	"github.com/illuminateeducation/rest-service-lib-go/pkg/validation"
	"net/http"
	"reflect"
	"strings"
)

// Parse through array type query parameters
func GetQueryParams(uri string, fb *db.FindBy, model interface{}, validator *validation.Validator) error {
	filters := strings.Split(uri, "?")

	if len(filters) > 1 {
		filters = strings.Split(filters[1], "&")

		for _, filter := range filters {
			equalIndex := strings.Index(filter, "=")

			// Array Params.
			if strings.Contains(filter, "[") {
				openBracketIndex := strings.Index(filter, "[")
				closeBracketIndex := strings.Index(filter, "]")
				key := filter[0:openBracketIndex]
				field := filter[openBracketIndex+1 : closeBracketIndex]

				// Checks if `=` was omitted from the parameter
				if equalIndex == -1 {
					return errors.New(key + ": '" + field + "' field cannot be blank.")
				}

				value := filter[equalIndex+1:]
				fields := db.GetJsonToDbMap(model)

				if _, ok := fields[field]; !ok {
					return errors.New(field + ": This property does not exist.")
				}

				// delete "search[field]" or "sort[field]"
				delete(fb.Conditions, filter[0:equalIndex])

				switch key {
				case "search":
					// Cannot search by id (UUID in the DB), by field of type DateTime, and by type bool
					if isFieldUUID(model, field) || isFieldDateTime(model, field) || isFieldBoolean(model, field) {
						return errors.New(key + ": cannot search by '" + field + "' field")
					}

					// Not allowed to search by an empty value.
					if value == "" {
						return errors.New(key + ": '" + field + "' field cannot be blank.")
					}

					fb.Search[field] = value
					break
				case "sort":
					// Sort direction can only be 'asc' or 'desc'.
					if value != "asc" && value != "desc" {
						return errors.New("The sort field '" + field + "' has an invalid value of '" + value + "'. The valid values are 'asc' or 'desc'.")
					}

					// Cannot sort by id
					if field == "id" {
						return errors.New(key + ": cannot sort by id")
					}

					// Cannot sort by password
					if field == "password" {
						return errors.New(key + ": cannot sort by" + field)
					}

					fb.OrderBy[field] = value
					break
				default:
					return errors.New(key + ": This property does not exist.")
					break
				}
			} else {
				// Non-array Params.

				// Checks if `=` was omitted from the parameter
				if equalIndex == -1 {
					return errors.New(filter + ": Cannot be blank.")
				}

				key := filter[0:equalIndex]
				fields := db.GetJsonToDbMap(model)
				value := filter[equalIndex+1:]

				// Ignore reserved words.
				if key == "sort" || key == "page" || key == "size" || key == "search" {
					continue
				}

				if _, ok := fields[key]; !ok {
					return errors.New(key + ": This property does not exist.")
				}

				// Check for UUIDs.
				if isFieldUUID(model, key) {
					err := ValidateId(value, validator)

					if err != nil {
						return errors.New(key + ": Invalid UUID v4.")
					}
				}

				// Not allowed to filter.
				if isFieldDateTime(model, key) {
					return errors.New(key + ": Not allowed to filter by.")
				}

				// Not allowed to filter by an empty value.
				if value == "" {
					return errors.New(key + ": Cannot be blank.")
				}
			}
		}
	}

	return nil
}

// WriteErrorResponse will construct and write a json encoded ErrorResponse to the Response Writer
func WriteErrorResponse(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response.NewErrorResponse(code, err.Error()))
}

// WriteBadRequestErrorResponse will construct and write a json encoded ErrorResponse to the Response Writer with a 400 error
func WriteBadRequestErrorResponse(w http.ResponseWriter, err error) {
	WriteErrorResponse(w, http.StatusBadRequest, err)
}

// Write404ErrorResponse will construct and write a json encoded ErrorResponse to the Response Writer with a 404 error
func Write404ErrorResponse(w http.ResponseWriter) {
	WriteErrorResponse(w, http.StatusNotFound, NotFound404)
}

// WriteSingleResponse will construct and write a json encoded SingleResponse to the Response Writer
func WriteSingleResponse(model interface{}, resourceType string, rm map[string]string, router *mux.Router, w http.ResponseWriter, r *http.Request, successfulStatusCode int) {
	sr, err := response.NewModelSingleResponse(model, rm, resourceType, router, r)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	resp, err := json.Marshal(sr)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(successfulStatusCode)
	w.Write(resp)
}

// FindModel will use the `id` variable from the url to attempt to find a model from the repository.  If none is found
// nil is returned.
func FindModel(model interface{}, rep db.Repository, r *http.Request) interface{} {
	vars := mux.Vars(r)
	err := rep.Find(model, vars["id"])
	if err != nil {
		return nil
	}

	return model
}

func GetRouteParams(model interface{}) map[string][]string {
	modelMap := db.StructToMap(model)
	routeParams := map[string][]string{
		route.CGET_ROUTE:   {},
		route.GET_ROUTE:    {"id", modelMap["id"].(string)},
		route.POST_ROUTE:   {},
		route.PATCH_ROUTE:  {"id", modelMap["id"].(string)},
		route.DELETE_ROUTE: {"id", modelMap["id"].(string)},
	}

	return routeParams
}

func ValidateId(id string, validator *validation.Validator) error {
	return (*validator).Var(id, "uuid4")
}

// cannot search by field of type UUID
func isFieldUUID(model interface{}, field string) bool {
	fields := structs.Fields(model)
	v := reflect.ValueOf(model)
	t := reflect.TypeOf(model)
	for i := range fields {
		validate := t.Field(i).Tag.Get("validate")
		if strings.Index(validate, "uuid") >= 0 && strings.ToLower(field) == strings.ToLower(v.Type().Field(i).Name) {
			return true
		}
	}

	return false
}

// cannot search by field of type types.NullDatetime
func isFieldDateTime(model interface{}, field string) bool {
	fields := structs.Fields(model)
	v := reflect.ValueOf(model)
	for i := range fields {
		t := v.Field(i).Type()
		if (t == reflect.TypeOf(types.NullDatetime{}) || t == reflect.TypeOf(types.Datetime{})) &&
			strings.ToLower(field) == strings.ToLower(v.Type().Field(i).Name) {
			return true
		}
	}

	return false
}

// cannot search by field of type bool
func isFieldBoolean(model interface{}, field string) bool {
	fields := structs.Fields(model)
	t := reflect.TypeOf(model)
	for i := range fields {
		if t.Field(i).Type.Name() == "bool" && strings.ToLower(field) == strings.ToLower(t.Field(i).Name) {
			return true
		}
	}

	return false
}
