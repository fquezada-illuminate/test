package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/illuminateeducation/rest-service-lib-go/pkg/response"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestNotFoundHandler_ServeHTTP(t *testing.T) {
	h := NotFoundHandler{}

	sr := strings.NewReader("")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/bad-route", sr)

	h.ServeHTTP(rr, req)

	if rr.Header().Get("Content-Type") != "application/json" {
		t.Error("Content-Type of 404 should be json")
	}

	if rr.Code != http.StatusNotFound {
		t.Error("Status code for not found should be 404")
	}

	e := response.ErrorResponse{}
	json.NewDecoder(rr.Body).Decode(&e)
	if e.Error.Message == "" || e.Error.Code == 0 {
		t.Error("Message and code cannot be empty")
	}
}

func TestMethodNotAllowedHandler_ServeHTTP(t *testing.T) {
	h := MethodNotAllowedHandler{}

	sr := strings.NewReader("")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", sr)

	h.ServeHTTP(rr, req)

	if rr.Header().Get("Content-Type") != "application/json" {
		t.Error("Content-Type of 415 should be json")
	}

	if rr.Code != http.StatusMethodNotAllowed {
		t.Error("Status code for Method Not Allowed should be 405")
	}

	e := response.ErrorResponse{}
	json.NewDecoder(rr.Body).Decode(&e)
	if e.Error.Message == "" || e.Error.Code == 0 {
		t.Error("Message and code cannot be empty")
	}
}

func TestCreateHealthCheckHandler(t *testing.T) {

	serviceId := "service-id"
	version := "api-version"
	releaseTag := "release-tag"
	os.Setenv("API_VERSION", version)
	os.Setenv("RELEASE_TAG", releaseTag)

	t.Run("Test passed handler function", func(t *testing.T) {
		passingFunc := func() error {
			return nil
		}
		rr := httptest.NewRecorder()
		sr := strings.NewReader("")
		req := httptest.NewRequest("GET", "/", sr)

		f := CreateHealthCheckHandler(passingFunc, serviceId)

		f.ServeHTTP(rr, req)

		// Check the response body is what we expect.
		expected := fmt.Sprintf(`{"status":"pass","serviceId":"%s","version":"%s","releaseId":"%s"}`, serviceId, version, releaseTag)
		if strings.Trim(rr.Body.String(), "\n") != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}

		if rr.Code != http.StatusOK {
			t.Errorf("when handler passes, the header should be 200 got %d", rr.Code)
		}
	})

	t.Run("Test failed handler function", func(t *testing.T) {
		err := "error message"
		failingFunc := func() error {
			return errors.New(err)
		}

		rr := httptest.NewRecorder()
		sr := strings.NewReader("")
		req := httptest.NewRequest("GET", "/", sr)

		f := CreateHealthCheckHandler(failingFunc, serviceId)

		f.ServeHTTP(rr, req)

		// Check the response body is what we expect.
		expected := fmt.Sprintf(`{"status":"fail","serviceId":"%s","version":"%s","releaseId":"%s","error":"%s"}`, serviceId, version, releaseTag, err)
		if strings.Trim(rr.Body.String(), "\n") != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}

		if rr.Code != http.StatusServiceUnavailable {
			t.Errorf("when handler passes, the header should be 503 got %d", rr.Code)
		}
	})
}

func TestGetSuccessfulBootHealthCheck(t *testing.T) {
	fn := GetSuccessfulBootHealthCheck()

	if reflect.TypeOf(fn).Kind() != reflect.Func {
		t.Errorf("GetSuccessfulBootHealthCheck should return a function")
	}

	if fn() != nil {
		t.Errorf("GetSuccessfulBootHealthCheck should return a function that always returns nil.")
	}

}
