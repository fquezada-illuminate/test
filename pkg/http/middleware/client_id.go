package middleware

import (
	"encoding/json"
	"github.com/illuminateeducation/rest-service-lib-go/pkg/response"
	"net/http"
)

func ClientId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-ied-client-id") == "" {
			// TODO: Use Error Response Helper/, whatever it becomes
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.NewErrorResponse(http.StatusBadRequest, http.StatusText(http.StatusBadRequest)+". HTTP Header X-Ied-Client-Id is missing or not a valid UUID format."))
			return
		}
		// TODO: Add UUID Validation

		next.ServeHTTP(w, r)
	})
}
