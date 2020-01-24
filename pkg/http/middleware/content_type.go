package middleware

import (
	"encoding/json"
	"github.com/illuminateeducation/rest-service-lib-go/pkg/response"
	"net/http"
)

func ContentType(accept string, respond string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", respond)

			// check content-type
			if r.Header.Get("Content-Type") != accept {
				// TODO: Use Error Response Helper/, whatever it becomes
				w.WriteHeader(http.StatusUnsupportedMediaType)
				json.NewEncoder(w).Encode(response.NewErrorResponse(http.StatusUnsupportedMediaType, http.StatusText(http.StatusUnsupportedMediaType)+". Expecting "+accept+" as Content-Type."))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func JsonContentType(next http.Handler) http.Handler {
	return ContentType("application/json", "application/json")(next)
}
