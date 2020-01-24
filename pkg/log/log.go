// log package contains ways to easily create most commonly used logs
package log

import (
	"net/http"
	"time"
)

// Request is for logging request structure
type Request struct {
	Time    time.Time
	Method  string
	Headers http.Header
	Client  string
	Url     string
}

// Response is for logging response structure
type Response struct {
	Headers     http.Header
	Body        string
	ServiceName string
	Time        time.Time
	Duration    time.Duration
}