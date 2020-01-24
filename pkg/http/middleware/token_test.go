package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestTokenReturnsError(t *testing.T) {
	os.Setenv("SERVICE_TOKEN", "invalid")
	r := httptest.NewRequest("GET", "/", nil)

	okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Next handler should not execute")
	})

	h := Token(okHandler)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected HTTP status code %d, got %d", http.StatusUnauthorized, w.Code)
	}

	if w.Body.String() == "" {
		t.Fatal("No response given.")
	}
}

func TestTokenPasses(t *testing.T) {
	os.Setenv("SERVICE_TOKEN", "valid")
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("x-ied-service-token", "valid")

	okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	h := Token(okHandler)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)

	if w.Body.String() != "ok" {
		t.Fatal("Invalid Service Token")
	}
}
