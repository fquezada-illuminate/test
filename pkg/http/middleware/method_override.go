package middleware

import (
	"net/http"
)

const HTTPMethodOverrideHeader = "X-HTTP-Method-Override"

func HandleHTTPMethodOverride(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			om := r.Header.Get(HTTPMethodOverrideHeader)
			if om == "PUT" || om == "PATCH" || om == "DELETE" {
				r.Method = om
			}
		}
		next.ServeHTTP(w, r)
	})
}