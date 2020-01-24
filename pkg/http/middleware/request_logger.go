package middleware

import (
	"bytes"
	"github.com/illuminateeducation/rest-service-lib-go/pkg/log"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"time"
)

// ResponseWriter overwrites the default writer so we can catch the response
type ResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

// Write writes the response
func (w ResponseWriter) Write(p []byte) (int, error) {
	return w.buf.Write(p)
}

// RequestLogger is a middleware logs the incoming request and outgoing response to a log
func RequestLogger(logger *logrus.Logger, serviceName string) func(next http.Handler) http.Handler {
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			timeStart := time.Now()

			requestLog := log.Request{
				timeStart.UTC(),
				r.Method,
				r.Header,
				r.Header.Get("x-ied-client-id"),
				r.URL.String(),
			}

			w2 := &ResponseWriter{w, &bytes.Buffer{}}

			defer func() {
				timeEnd := time.Now()
				responseLog := log.Response{
					w.Header(),
					w2.buf.String(),
					serviceName,
					timeEnd.UTC(),
					timeEnd.Sub(timeStart),
				}

				logger.
					WithField("request", requestLog).
					WithField("response", responseLog).
					Infof("%v API Call", serviceName)

				io.Copy(w, w2.buf)
			}()

			next.ServeHTTP(w2, r)
		})
	}
}