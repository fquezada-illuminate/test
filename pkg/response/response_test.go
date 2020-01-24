package response

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/illuminateeducation/rest-service-lib-go/pkg/http/route"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const RESOURCE_TYPE = "model"

type Model struct {
	Id string `json:"id"`
}

var routeNames = map[string]string{
	route.CGET_ROUTE:   route.CGET_ROUTE,
	route.GET_ROUTE:    route.GET_ROUTE,
	route.POST_ROUTE:   route.POST_ROUTE,
	route.PATCH_ROUTE:  route.PATCH_ROUTE,
	route.DELETE_ROUTE: route.DELETE_ROUTE,
}

func TestCreateSingleResponse(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	router := mux.NewRouter()

	t.Run("Valid instantiation of SingleResponse", func(t *testing.T) {
		sr, err := CreateSingleResponse(struct{}{}, "type", router, req)

		if err != nil {
			t.Errorf("Error was not expected when creating SingleResponse: %s", err.Error())
		}

		if sr.router == nil {
			t.Error("Route was not set when creating the SingleResponse")
		}

		if sr.request == nil {
			t.Error("Request was not set when creating the SingleResponse")
		}
	})

	t.Run("Provide an invalid RESOURCE_TYPE", func(t *testing.T) {
		_, err := CreateSingleResponse(struct{}{}, "", router, req)

		if err == nil {
			t.Error("Error was not expected when creating SingleResponse with empty RESOURCE_TYPE")
		}
	})
}

func TestSingleResponse_AddLink(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", strings.NewReader("body"))
	router := mux.NewRouter()
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {

	}).Name("endpoint")
	router.HandleFunc("/{id}/test", func(writer http.ResponseWriter, request *http.Request) {

	}).Name("identifier")

	sr, _ := CreateSingleResponse(struct{}{}, "test-resource", router, req)

	t.Run("Return an error for bad route", func(t *testing.T) {
		routeParams := make(map[string]string)
		err := sr.AddLink("get", "Relation", "badRoute", routeParams)
		if err == nil {
			t.Error("If a route doesn't exist it should return an error")
		}
	})

	t.Run("Return an error when can't generate url", func(t *testing.T) {
		routeParams := make(map[string]string)
		err := sr.AddLink("get", "Relation", "endpoints", routeParams)
		if err == nil {
			t.Error("If there's a problem generating the url, an error should be returned. Endpoint doesn't exist")
		}
	})

	t.Run("Return an error when can't generate url", func(t *testing.T) {
		routeParams := make(map[string]string)
		err := sr.AddLink("get", "Identifier", "identifier", routeParams)
		if err == nil {
			t.Error("If there's a problem generating the url, an error should be returned. Param doesn't exist")
		}
	})

	t.Run("AddLink should create and be added to the list of links", func(t *testing.T) {
		routeParams := map[string]string{
			"id": "1",
		}
		err := sr.AddLink("get", "Relation", "endpoint", routeParams)
		if err != nil {
			t.Error("A valid route should not return an error")
		}

		links := sr.links
		if len(links) != 1 {
			t.Error("A link was not added to the list")
		}

		err = sr.AddLink("get", "Identifier", "identifier", routeParams)
		if err != nil {
			t.Error("A valid route should not return an error")
		}

		links = sr.links
		if len(links) != 2 {
			t.Error("A link with a parameter was not added to the list")
		}
	})
}

func TestSingleResponse_MarshalJSON(t *testing.T) {

	req := httptest.NewRequest("GET", "http://example.com", nil)
	router := mux.NewRouter()
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {

	}).Name("endpoint")

	t.Run("Response output should be unmarshallable from json", func(t *testing.T) {
		type TestModel struct {
			Name string
		}
		r, err := CreateSingleResponse(TestModel{Name: "test-name"}, "type", router, req)
		if err != nil {
			t.Error("Error was not expected when creating SingleResponse")
		}
		routeParams := make(map[string]string)
		r.AddLink("GET", "self", "endpoint", routeParams)

		b, err := json.Marshal(r)
		if err != nil {
			t.Errorf("An error was returned from GetResponse when not expected: %s", err.Error())
		}

		testStruct := struct {
			Model TestModel `json:"type"`
			Links []Link    `json:"links"`
		}{}

		err = json.NewDecoder(bytes.NewReader(b)).Decode(&testStruct)
		if err != nil {
			t.Errorf("Failed to recode the Response to json: %s", err.Error())
		}
	})

	t.Run("Bad marshalling should return an error", func(t *testing.T) {
		r, err := CreateSingleResponse(make(chan int), "type", router, req)
		if err != nil {
			t.Error("Error was not expected when creating a SingleResponse")
		}

		_, err = json.Marshal(r)
		if err == nil {
			t.Error("A Response should return an error if it can't be unMarshalled")
		}
	})
}

func TestCreateCollectionResponse(t *testing.T) {
	cm := CollectionMetadata{}
	req := httptest.NewRequest("GET", "http://example.com", nil)
	router := mux.NewRouter()

	t.Run("Valid instantiation of CollectionResponse", func(t *testing.T) {
		cr, err := CreateCollectionResponse(cm, "type", router, req)

		if err != nil {
			t.Error("Error was not expected when creating CollectionResponse")
		}

		if cr.router == nil {
			t.Error("Route was not set when creating the CollectionResponse")
		}

		if cr.request == nil {
			t.Error("Request was not set when creating the CollectionResponse")
		}
	})

	t.Run("Provide an invalid RESOURCE_TYPE", func(t *testing.T) {
		_, err := CreateCollectionResponse(cm, "", router, req)

		if err == nil {
			t.Error("Error was not expected when creating CollectionResponse")
		}
	})
}

func TestCollectionResponse_AddItem(t *testing.T) {
	cm := CollectionMetadata{}
	req := httptest.NewRequest("GET", "http://example.com", nil)
	router := mux.NewRouter()

	t.Run("Add a new Item", func(t *testing.T) {
		cr, _ := CreateCollectionResponse(cm, "type", router, req)
		cr.AddItem(SingleResponse{})

		if len(cr.items) != 1 {
			t.Errorf("Expected one item in the list, got %v", len(cr.items))
		}
	})
}

func TestCollectionResponse_MarshalJSON(t *testing.T) {
	cm := CollectionMetadata{}
	req := httptest.NewRequest("GET", "http://example.com", nil)
	router := mux.NewRouter()
	router.HandleFunc("", func(writer http.ResponseWriter, request *http.Request) {

	}).Name("endpoint")

	t.Run("Return a full marshalled CollectionResponse", func(t *testing.T) {
		cr, _ := CreateCollectionResponse(cm, "type", router, req)
		routeParams := make(map[string]string)
		cr.AddLink("get", "self", "endpoint", routeParams)

		b, _ := json.Marshal(cr)
		ret := make(map[string]interface{})
		json.Unmarshal(b, &ret)

		if _, ok := ret["items"]; !ok {
			t.Error("items key was expected to be returned")
		}

		if _, ok := ret["metadata"]; !ok {
			t.Error("metadata key was expected to be returned")
		}

		if _, ok := ret["links"]; !ok {
			t.Error("links key was expected to be returned")
		}
	})

	t.Run("If there are no link the key shouldn't be present", func(t *testing.T) {
		cr, _ := CreateCollectionResponse(cm, "type", router, req)

		b, _ := json.Marshal(cr)
		ret := make(map[string]interface{})
		json.Unmarshal(b, &ret)

		if _, ok := ret["links"]; ok {
			t.Error("links key should not be present")
		}
	})

}

func TestNewModelSingleResponse(t *testing.T) {
	model := Model{Id: "test-id"}

	t.Run("Return a valid single response", func(t *testing.T) {
		router := mux.NewRouter()
		router.Path("/models").Name(route.CGET_ROUTE)
		router.Path("/models/id").Name(route.GET_ROUTE)
		router.Path("/models/id").Name(route.PATCH_ROUTE)
		router.Path("/models/id").Name(route.DELETE_ROUTE)

		req, _ := http.NewRequest("GET", "http://example.com", nil)
		sr, err := NewModelSingleResponse(model, routeNames, RESOURCE_TYPE, router, req)

		if err != nil {
			t.Errorf("Got an error when one wasn't expected: %s", err.Error())
		}

		if len(sr.GetLinks()) != 4 {
			t.Errorf("Expected 4 links, got: %d", len(sr.GetLinks()))
		}
	})

	t.Run("Return errors when paths don't exist", func(t *testing.T) {
		router := mux.NewRouter()

		req, _ := http.NewRequest("GET", "http://example.com", nil)
		_, err := NewModelSingleResponse(model, routeNames, RESOURCE_TYPE, router, req)

		if err == nil {
			t.Errorf("An error was expected but none was returned")
		}
	})
}

func TestNewModelCollectionResponse(t *testing.T) {

	cm := CollectionMetadata{}

	t.Run("Return a valid SingleResponse", func(t *testing.T) {
		router := mux.NewRouter()
		router.Path("/").Name(routeNames[route.CGET_ROUTE])
		router.Path("/").Name(routeNames[route.POST_ROUTE])
		req, _ := http.NewRequest("GET", "http://example.com", nil)

		cr, err := NewModelCollectionResponse(cm, routeNames, RESOURCE_TYPE, router, req)

		if err != nil {
			t.Errorf("Got an error when one wasn't expected: %s", err.Error())
		}

		if len(cr.GetLinks()) != 4 {
			t.Errorf("Expected 4 links got: %d", len(cr.GetLinks()))
		}
	})

	t.Run("Return errors when paths don't exist", func(t *testing.T) {
		router := mux.NewRouter()
		req, _ := http.NewRequest("GET", "http://example.com", nil)
		_, err := NewModelCollectionResponse(cm, routeNames, RESOURCE_TYPE, router, req)

		if err == nil {
			t.Errorf("An error was expected but none was returned")
		}
	})
}
