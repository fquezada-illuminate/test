package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleHTTPMethodOverride(t *testing.T) {

	tests := []struct {
		httpMethod     string
		overrideMethod string
		expectedMethod string
	}{
		{
			httpMethod: "GET",
			overrideMethod: "PATCH",
			expectedMethod: "GET",
		},
		{
			httpMethod: "POST",
			overrideMethod: "PATCH",
			expectedMethod: "PATCH",
		},
		{
			httpMethod: "POST",
			overrideMethod: "PUT",
			expectedMethod: "PUT",
		},
		{
			httpMethod: "POST",
			overrideMethod: "DELETE",
			expectedMethod: "DELETE",
		},
		{
			httpMethod: "POST",
			overrideMethod: "GET",
			expectedMethod: "POST",
		},
		{
			httpMethod: "PATCH",
			overrideMethod: "DELETE",
			expectedMethod: "PATCH",
		},
		{
			httpMethod: "DELETE",
			overrideMethod: "PATCH",
			expectedMethod: "DELETE",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Test %s with %s override", test.httpMethod, test.overrideMethod), func(t *testing.T) {
			okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != test.expectedMethod {
					t.Errorf("Unexpected HTTP Method.\n Expected:\t %v \n Got:\t %v", test.expectedMethod, r.Method)
				}
			})

			r := httptest.NewRequest(test.httpMethod, "/", nil)
			r.Header.Set(HTTPMethodOverrideHeader, test.overrideMethod)
			w := httptest.NewRecorder()
			h := HandleHTTPMethodOverride(okHandler)
			h.ServeHTTP(w, r)
		})
	}
}
