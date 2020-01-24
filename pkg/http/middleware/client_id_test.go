package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientIdReturnsError(t *testing.T) {
	// Set up server
	okHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		t.Fatal("Next handler should not execute.")
	})
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h := ClientId(okHandler)
	h.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected HTTP status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	if w.Body.String() == "" {
		t.Fatal("No response given.")
	}

	//if w.Header().Get("content-type") != "application/json" {
	//	t.Errorf("Expected %s got %s", "application/json", w.Header().Get("content-type"))
	//}
	//
	//errResponse := &ErrorResponse{}
	//err := json.Unmarshal(w.Body.Bytes(), errResponse)
	//
	//if err != nil {
	//	t.Fatal("Unabled to Unmsarshal")
	//}
	//
	//if errResponse.Code != http.StatusBadRequest {
	//	t.Errorf("Expected error body code to be %d got %d", http.StatusBadRequest, errResponse.Code)
	//}
	//
	//if errResponse.Message != http.StatusText(http.StatusBadRequest) {
	//	t.Errorf("Expected error message %s got %s", http.StatusText(http.StatusBadRequest), errResponse.Message)
	//}
}

func TestClientIdPasses(t *testing.T) {
	// Set up server
	okHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("ok"))
	})
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("x-ied-client-id", "1")
	w := httptest.NewRecorder()
	h := ClientId(okHandler)
	h.ServeHTTP(w, r)

	if w.Body.String() != "ok" {
		t.Fatal("x-ied-client-id was not set")
	}
}
