package middleware

import (
	"encoding/json"
	"github.com/illuminateeducation/rest-service-lib-go/pkg/response"
	"net/http"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				w.Header().Set("content-type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response.NewErrorResponse(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
