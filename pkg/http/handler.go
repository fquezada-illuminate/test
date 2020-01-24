package http

import (
	"encoding/json"
	"fmt"
	"github.com/illuminateeducation/rest-service-lib-go/pkg/response"
	"net/http"
	"os"
)

type NotFoundHandler struct {
}

func (h NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	method := r.Method
	uri := r.RequestURI
	errorMessage := fmt.Sprintf(`No route found for "%s %s"`, method, uri)
	json.NewEncoder(w).Encode(response.NewErrorResponse(http.StatusNotFound, errorMessage))
}

type MethodNotAllowedHandler struct {
}

func (h MethodNotAllowedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(response.NewErrorResponse(http.StatusMethodNotAllowed, "Method not allowed"))
}

type HealthCheckFunc func() error

//CreateHealthCheckHandler is a generic function that will produce the health check endpoint.
//The healthCheckFunc will be used to determine if the application is booted up properly.
func CreateHealthCheckHandler(healthCheckFunc HealthCheckFunc, serviceId string) http.HandlerFunc {
	status := "pass"
	statusCode := http.StatusOK
	err := healthCheckFunc()
	errorMessage := ""
	if err != nil {
		status = "fail"
		statusCode = http.StatusServiceUnavailable
		errorMessage = err.Error()
	}

	ret := struct {
		Status    string `json:"status"`
		ServiceId string `json:"serviceId"`
		Version   string `json:"version"`
		ReleaseId string `json:"releaseId"`
		Error     string `json:"error,omitempty"`
	}{
		Status:    status,
		ServiceId: serviceId,
		Version:   os.Getenv("API_VERSION"),
		ReleaseId: os.Getenv("RELEASE_TAG"),
		Error:     errorMessage,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(ret)
	}
}

//GetSuccessfulBootHealthCheck will return a function that always returns nil.  There will be services that will
//panic if the binary can't boot up, this function will be used when that is the case.
func GetSuccessfulBootHealthCheck() HealthCheckFunc {
	return func() error {
		return nil
	}
}
