package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContentTypeReturnsError(t *testing.T) {
	// Set up server
	okHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		t.Fatal("Next handler should not execute.")
	})
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	cType := "application/json"
	h := ContentType(cType, cType)(okHandler)
	h.ServeHTTP(w, r)

	if w.Code != http.StatusUnsupportedMediaType {
		t.Errorf("Expected HTTP Status Code %d got %d", http.StatusUnsupportedMediaType, w.Code)
	}

	if w.Header().Get("Content-Type") != cType {
		t.Errorf("Expected %s Content-Type, got %s", cType, w.Header().Get("Content-Type"))
	}

	if w.Body.String() == "" {
		t.Fatal("No response given.")
	}
}

func TestContentTypePasses(t *testing.T) {
	// Set up server
	okHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("ok"))
	})
	cType := "application/json"
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Content-Type", cType)
	w := httptest.NewRecorder()
	h := ContentType(cType, cType)(okHandler)
	h.ServeHTTP(w, r)

	if w.Body.String() != "ok" {
		t.Fatal("Content-Type is invalid")
	}
}

func TestJsonContentTypeError(t *testing.T) {
	okHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		t.Fatal("Next handler should not execute.")
	})
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	cType := "application/json"
	h := JsonContentType(okHandler)

	h.ServeHTTP(w, r)

	if w.Code != http.StatusUnsupportedMediaType {
		t.Errorf("Expected HTTP Status Code %d got %d", http.StatusUnsupportedMediaType, w.Code)
	}

	if w.Header().Get("Content-Type") != cType {
		t.Errorf("Expected %s Content-Type, got %s", cType, w.Header().Get("Content-Type"))
	}

	if w.Body.String() == "" {
		t.Fatal("No response given.")
	}
}

func TestJsonContentTypePasses(t *testing.T) {
	okHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("ok"))
	})

	cType := "application/json"
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Content-Type", cType)
	w := httptest.NewRecorder()

	h := JsonContentType(okHandler)

	h.ServeHTTP(w, r)

	if w.Header().Get("Content-Type") != cType {
		t.Errorf("Expected %s Content-Type, got %s", cType, w.Header().Get("Content-Type"))
	}

	if w.Body.String() != "ok" {
		t.Fatal("Content-Type is invalid" + w.Body.String())
	}
}
