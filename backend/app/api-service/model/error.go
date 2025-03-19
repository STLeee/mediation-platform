package model

import (
	"net/http"
)

// HttpStatusCodeError struct for HTTP status code error
type HttpStatusCodeError struct {
	StatusCode int
	Message    string
	Err        error
}

// Error returns the error message
func (e HttpStatusCodeError) Error() string {
	message := e.Message
	if message == "" {
		message = http.StatusText(e.StatusCode)
	}
	return message
}

// Unwrap returns the wrapped error
func (e HttpStatusCodeError) Unwrap() error {
	return e.Err
}
