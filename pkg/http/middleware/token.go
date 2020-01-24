package middleware

import (
	"encoding/json"
	"github.com/illuminateeducation/rest-service-lib-go/pkg/response"
	"net/http"
	"os"
)

func Token(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-ied-service-token") != os.Getenv("SERVICE_TOKEN") {
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response.NewErrorResponse(http.StatusUnauthorized, http.StatusText(http.StatusBadRequest)+". Invalid Token in Headers."))
			return
		}
		next.ServeHTTP(w, r)
	})
}
