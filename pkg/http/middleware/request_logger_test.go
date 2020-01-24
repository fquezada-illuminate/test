package middleware

import (
	"bytes"
	"github.com/illuminateeducation/rest-service-lib-go/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRequestLogger(t *testing.T) {
	logger, hook := test.NewNullLogger()
	service := "Service Name"

	requestLogger := RequestLogger(logger, service)

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	okHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if _, ok := w.(*ResponseWriter); !ok {
			t.Errorf("invalid middleware chaining. Expected %s writer, got %s", reflect.TypeOf(&ResponseWriter{}), reflect.TypeOf(w))
		}
	})

	h := requestLogger(okHandler)
	h.ServeHTTP(w, r)

	lastEntry := hook.LastEntry()

	if _, ok := logger.Formatter.(*logrus.JSONFormatter); !ok {
		t.Errorf("invalid log format. Expected %s got %s", reflect.TypeOf(&logrus.JSONFormatter{}), reflect.TypeOf(logger.Formatter))
	}

	if lastEntry.Level != logrus.InfoLevel {
		t.Errorf("invalid log level. Expected %s got %s", logrus.InfoLevel, lastEntry.Level)
	}

	logData := hook.LastEntry().Data

	if _, ok := logData["request"]; !ok {
		t.Errorf("invalid log format. Log does not have request field")
	}

	if _, ok := logData["response"]; !ok {
		t.Errorf("invalid log format. Log does not have response field")
	}

	if _, ok := logData["request"].(log.Request); !ok {
		t.Errorf("invalid log request struct type. Expected %s got %s", reflect.TypeOf(log.Request{}), reflect.TypeOf(logData["request"]))
	}

	if _, ok := logData["response"].(log.Response); !ok {
		t.Errorf("invalid log request struct type. Expected %s got %s", reflect.TypeOf(log.Request{}), reflect.TypeOf(logData["response"]))
	}

	if lastEntry.Message != service + " API Call" {
		t.Errorf("invalid log message. Expected %s got %s", service + " API Call", lastEntry.Message)
	}
}

func TestResponseWriter_Write(t *testing.T) {
	s := "Test"

	w := httptest.NewRecorder()
	rw := ResponseWriter{w, &bytes.Buffer{}}
	i, err := rw.Write([]byte(s))

	if err != nil {
		t.Error("Unexpected write error")
	}
	
	if i != len(s) {
		t.Error("Unexpected write length")
	}

}
