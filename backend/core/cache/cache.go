package cache

import "strings"

type CacheErrorType string

const (
	CacheErrorTypeServerError    CacheErrorType = "server_error"
	CacheErrorTypeOperationError CacheErrorType = "operation_error"
)

var CacheErrorDefaultMessages = map[CacheErrorType]string{
	CacheErrorTypeServerError:    "server error",
	CacheErrorTypeOperationError: "operation error",
}

// CacheError struct for database error
type CacheError struct {
	ErrType CacheErrorType
	Message string
	Err     error
}

// Error returns the error message
func (e CacheError) Error() string {
	message := e.Message
	if message == "" {
		if defaultMessage, ok := CacheErrorDefaultMessages[e.ErrType]; ok {
			message = defaultMessage
		}
	}
	if e.Err != nil {
		message = strings.Join([]string{message, e.Err.Error()}, ": ")
	}
	return message
}

// Unwrap returns the wrapped error
func (e CacheError) Unwrap() error {
	return e.Err
}
