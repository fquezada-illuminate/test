package svc

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/illuminateeducation/rest-service-lib-go/pkg/db"
	"github.com/illuminateeducation/rest-service-lib-go/pkg/http/route"
	"github.com/illuminateeducation/rest-service-lib-go/pkg/types"
	"github.com/illuminateeducation/rest-service-lib-go/pkg/validation"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type Model struct {
	Id          string             `json:"id" db:"id" structs:"id" validate:"uuid4"`
	Name        types.NullString   `json:"name" db:"name" structs:"name,omitnested" validate:"min=3,max=63"`
	CreatedAt   types.NullDatetime `json:"createdAt" db:"created_at" structs:"created_at,omitnested"`
	ShowProduct bool               `json:"showProduct" db:"show_product" structs:"show_product"`
}

type mockRepo struct {
	db.BaseRepository
}

func (r mockRepo) Find(object interface{}, id string) error {
	if err := r.IsPointer(object); err != nil {
		return err
	}

	obj, ok := object.(*Model)

	if !ok {
		log.Fatal("Error casting interface{} to *svc.Model")
	}

	obj.Id = "c24b2909-92e3-4266-ac13-95ac9f24388f"
	obj.Name = types.NullString{}
	obj.CreatedAt = types.NullDatetime{}

	return nil
}

func TestNoFilters(t *testing.T) {
	fb := &db.FindBy{}
	validator := validation.Singleton()
	model := Model{}

	params := GetQueryParams("/", fb, model, validator)

	if params != nil {
		t.Errorf("Expected nil, got %v", params)
	}
}

func TestEqualSignMissing(t *testing.T) {
	fb := &db.FindBy{}
	validator := validation.Singleton()
	model := Model{}

	params := GetQueryParams("?name", fb, model, validator)

	if _, ok := params.(error); !ok {
		t.Errorf("Expected error, got %v", params)
	}

	if params.Error() != "name: Cannot be blank." {
		t.Errorf("Expected \"name: cannot be blank\", got \"%s\"", params.Error())
	}
}

func TestIgnoreReservedWords(t *testing.T) {
	fb := &db.FindBy{}
	validator := validation.Singleton()
	model := Model{}

	params := GetQueryParams("?page=7", fb, model, validator)

	if params != nil {
		t.Errorf("Expected nil, got %v", params)
	}
}

func TestPropertyDoesNotExist(t *testing.T) {
	fb := &db.FindBy{}
	validator := validation.Singleton()
	model := Model{}

	params := GetQueryParams("?test=test", fb, model, validator)

	if _, ok := params.(error); !ok {
		t.Errorf("Expected error, got %v", params)
	}

	if params.Error() != "test: This property does not exist." {
		t.Errorf("Expected \"test: This property does not exist.\", got \"%s\"", params.Error())
	}
}

func TestInvalidUuid(t *testing.T) {
	fb := &db.FindBy{}
	validator := validation.Singleton()
	model := Model{}

	params := GetQueryParams("?id=test", fb, model, validator)

	if _, ok := params.(error); !ok {
		t.Errorf("Expected error, got %v", params)
	}

	if params.Error() != "id: Invalid UUID v4." {
		t.Errorf("Expected \"id: Invalid UUID v4.\", got \"%s\"", params.Error())
	}
}

func TestCannotFilterByDatetime(t *testing.T) {
	fb := &db.FindBy{}
	validator := validation.Singleton()
	model := Model{}

	params := GetQueryParams("?createdAt=test", fb, model, validator)

	if _, ok := params.(error); !ok {
		t.Errorf("Expected error, got %v", params)
	}

	if params.Error() != "createdAt: Not allowed to filter by." {
		t.Errorf("Expected \"createdAt: Not allowed to filter by.\", got \"%s\"", params.Error())
	}
}

func TestFilterValueCannotBeBlank(t *testing.T) {
	fb := &db.FindBy{}
	validator := validation.Singleton()
	model := Model{}

	params := GetQueryParams("?name=", fb, model, validator)

	if _, ok := params.(error); !ok {
		t.Errorf("Expected error, got %v", params)
	}

	if params.Error() != "name: Cannot be blank." {
		t.Errorf("Expected \"name: cannot be blank\", got \"%s\"", params.Error())
	}
}

func TestCannotSearchByUUID(t *testing.T) {
	fb := &db.FindBy{}
	validator := validation.Singleton()
	model := Model{}

	params := GetQueryParams("?search[id]=test", fb, model, validator)

	if _, ok := params.(error); !ok {
		t.Errorf("Expected error, got %v", params)
	}

	if params.Error() != "search: cannot search by 'id' field" {
		t.Errorf("Expected \"search: cannot search by 'id' field\", got \"%s\"", params.Error())
	}
}

func TestCannotSearchByBoolean(t *testing.T) {
	fb := &db.FindBy{}
	validator := validation.Singleton()
	model := Model{}

	params := GetQueryParams("?search[showProduct]=true", fb, model, validator)

	if _, ok := params.(error); !ok {
		t.Errorf("Expected error, got %v", params)
	}

	if params.Error() != "search: cannot search by 'showProduct' field" {
		t.Errorf("Expected \"search: cannot search by 'id' field\", got \"%s\"", params.Error())
	}
}

func TestSearchCannotBeBlank(t *testing.T) {
	fb := &db.FindBy{}
	validator := validation.Singleton()
	model := Model{}

	params := GetQueryParams("?search[name]=", fb, model, validator)

	if _, ok := params.(error); !ok {
		t.Errorf("Expected error, got %v", params)
	}

	if params.Error() != "search: 'name' field cannot be blank." {
		t.Errorf("Expected \"search: 'name' field cannot be blank.\", got \"%s\"", params.Error())
	}
}

func TestInvalidSortDirection(t *testing.T) {
	fb := &db.FindBy{}
	validator := validation.Singleton()
	model := Model{}

	params := GetQueryParams("?sort[name]=test", fb, model, validator)

	if _, ok := params.(error); !ok {
		t.Errorf("Expected error, got %v", params)
	}

	if params.Error() != "The sort field 'name' has an invalid value of 'test'. The valid values are 'asc' or 'desc'." {
		t.Errorf("Expected \"The sort field 'name' has an invalid value of 'test'. The valid values are 'asc' or 'desc'.\", got \"%s\"", params.Error())
	}
}

func TestSortById(t *testing.T) {
	fb := &db.FindBy{}
	validator := validation.Singleton()
	model := Model{}

	params := GetQueryParams("?sort[id]=asc", fb, model, validator)

	if _, ok := params.(error); !ok {
		t.Errorf("Expected error, got %v", params)
	}

	if params.Error() != "sort: cannot sort by id" {
		t.Errorf("Expected \"sort: cannot sort by id\", got \"%s\"", params.Error())
	}
}

func TestInvalidArrayType(t *testing.T) {
	fb := &db.FindBy{}
	validator := validation.Singleton()
	model := Model{}

	params := GetQueryParams("?test[name]=test", fb, model, validator)

	if _, ok := params.(error); !ok {
		t.Errorf("Expected error, got %v", params)
	}

	if params.Error() != "test: This property does not exist." {
		t.Errorf("Expected \"test: This property does not exist.\", got \"%s\"", params.Error())
	}
}

func TestEqualSignMissingFromArray(t *testing.T) {
	fb := &db.FindBy{}
	validator := validation.Singleton()
	model := Model{}

	params := GetQueryParams("?search[name]", fb, model, validator)

	if _, ok := params.(error); !ok {
		t.Errorf("Expected error, got %v", params)
	}

	if params.Error() != "search: 'name' field cannot be blank." {
		t.Errorf("Expected \"search: 'name' field cannot be blank.\", got \"%s\"", params.Error())
	}
}

func TestInvalidArrayKey(t *testing.T) {
	fb := &db.FindBy{}
	validator := validation.Singleton()
	model := Model{}

	params := GetQueryParams("?search[test]=test", fb, model, validator)

	if _, ok := params.(error); !ok {
		t.Errorf("Expected error, got %v", params)
	}

	if params.Error() != "test: This property does not exist." {
		t.Errorf("Expected \"test: This property does not exist.\", got \"%s\"", params.Error())
	}
}

func TestValidSearch(t *testing.T) {
	fb := &db.FindBy{
		Search: make(map[string]interface{}),
	}

	validator := validation.Singleton()
	model := Model{}
	params := GetQueryParams("?search[name]=test", fb, model, validator)

	if params != nil {
		t.Errorf("Expected nil, got %v", params)
	}

	expected := map[string]interface{}{
		"name": "test",
	}

	if !reflect.DeepEqual(fb.Search, expected) {
		t.Errorf("Expected %v, got %v", expected, fb.Search)
	}
}

func TestValidSort(t *testing.T) {
	fb := &db.FindBy{
		OrderBy: make(map[string]interface{}),
	}

	validator := validation.Singleton()
	model := Model{}
	params := GetQueryParams("?sort[name]=asc", fb, model, validator)

	if params != nil {
		t.Errorf("Expected nil, got %v", params)
	}

	expected := map[string]interface{}{
		"name": "asc",
	}

	if !reflect.DeepEqual(fb.OrderBy, expected) {
		t.Errorf("Expected %v, got %v", expected, fb.OrderBy)
	}
}

func TestWriteBadRequestErrorResponse(t *testing.T) {
	w := httptest.NewRecorder()
	WriteBadRequestErrorResponse(w, errors.New("error, error"))

	if w.Code != 400 {
		t.Errorf("Expected status code %d, got %d", 400, w.Code)
	}

	equal, err := IsEqualJson(w.Body.String(), "{\"error\":{\"code\":400,\"message\":\"error, error\"}}")

	if err != nil {
		t.Error(err.Error())
	}

	if !equal {
		t.Errorf("Expected \"{\"error\":{\"code\":400,\"message\":\"error, error\"}}\", got %v", w.Body.String())
	}
}

func TestWrite404ErrorResponse(t *testing.T) {
	w := httptest.NewRecorder()
	Write404ErrorResponse(w)

	if w.Code != 404 {
		t.Errorf("Expected status code %d, got %d", 404, w.Code)
	}

	equal, err := IsEqualJson(w.Body.String(), "{\"error\":{\"code\":404,\"message\":\"not found\"}}")

	if err != nil {
		t.Error(err.Error())
	}

	if !equal {
		t.Errorf("Expected \"{\"error\":{\"code\":404,\"message\":\"not found\"}}\", got \"%v\"", w.Body.String())
	}
}

func TestGetRouteParams(t *testing.T) {
	model := Model{}
	params := GetRouteParams(model)

	expected := map[string][]string{
		route.CGET_ROUTE:   {},
		route.GET_ROUTE:    {"id", ""},
		route.POST_ROUTE:   {},
		route.PATCH_ROUTE:  {"id", ""},
		route.DELETE_ROUTE: {"id", ""},
	}

	if !reflect.DeepEqual(params, expected) {
		t.Errorf("Expected %v, got %v", expected, params)
	}
}

func TestFindModel(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	model := FindModel(&Model{}, mockRepo{}, req)
	obj, ok := model.(*Model)

	if !ok {
		log.Fatal("Error casting interface{} to *svc.Model")
	}

	if obj.Id != "c24b2909-92e3-4266-ac13-95ac9f24388f" {
		t.Errorf("Expected the model Id to be %s, got %v", "c24b2909-92e3-4266-ac13-95ac9f24388f", obj.Id)
	}
}

func TestFindModelError(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	model := FindModel(Model{}, mockRepo{}, req)

	if model != nil {
		t.Errorf("Expected nil, got %v", model)
	}
}

func TestWriteSingleResponse(t *testing.T) {
	router := mux.NewRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	WriteSingleResponse(Model{"c24b2909-92e3-4266-ac13-95ac9f24388f", types.NullString{}, types.NullDatetime{}, true}, "model", make(map[string]string), router, w, req, 200)
	equal, err := IsEqualJson(w.Body.String(), "{\"model\":{\"id\":\"c24b2909-92e3-4266-ac13-95ac9f24388f\",\"name\":null,\"createdAt\":null,\"showProduct\":true}}")

	if err != nil {
		t.Error(err.Error())
	}

	if !equal {
		t.Errorf("Expected \"{\"model\":{\"id\":\"c24b2909-92e3-4266-ac13-95ac9f24388f\",\"name\":null,\"createdAt\":null,\"showProduct\":true}}\", got \"%v\"", w.Body.String())
	}
}

func TestWriteSingleResponseError(t *testing.T) {
	router := mux.NewRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	rm := make(map[string]string)
	rm[route.CGET_ROUTE] = "test"
	WriteSingleResponse("", "model", rm, router, w, req, 200)
	equal, err := IsEqualJson(w.Body.String(), "{\"error\":{\"code\":400,\"message\":\"Route test does not exist\"}}")

	if err != nil {
		t.Error(err.Error())
	}

	if !equal {
		t.Errorf("Expected \"{\"error\":{\"code\":400,\"message\":\"Route test does not exist\"}}\", got \"%v\"", w.Body.String())
	}
}

func TestWriteSingleResponseUnmarshalError(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/test", Handler).Name("cget_instance")
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	rm := make(map[string]string)
	rm[route.CGET_ROUTE] = "cget_instance"
	WriteSingleResponse(make(chan int), "model", rm, router, w, req, 200)

	equal, err := IsEqualJson(w.Body.String(), "{\"error\":{\"code\":400,\"message\":\"json: error calling MarshalJSON for type response.SingleResponse: json: unsupported type: chan int\"}}")

	if err != nil {
		t.Error(err.Error())
	}

	if !equal {
		t.Errorf("Expected \"{\"error\":{\"code\":400,\"message\":\"json: error calling MarshalJSON for type response.SingleResponse: json: unsupported type: chan int\"}}\", got \"%v\"", w.Body.String())
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}

func IsEqualJson(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	err := json.Unmarshal([]byte(s1), &o1)

	if err != nil {
		return false, err
	}

	err = json.Unmarshal([]byte(s2), &o2)

	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(o1, o2), nil
}
