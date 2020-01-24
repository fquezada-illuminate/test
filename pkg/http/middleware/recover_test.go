package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecoverReturnsError(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)

	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Error")
	})

	h := Recover(panicHandler)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected HTTP status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
