package db

import "strings"

type DBErrorType string

const (
	DBErrorTypeServerError DBErrorType = "server_error"
	DBErrorConfigError     DBErrorType = "config_error"
)

var DBErrorDefaultMessages = map[DBErrorType]string{
	DBErrorTypeServerError: "server error",
	DBErrorConfigError:     "config error",
}

// DBError struct for database error
type DBError struct {
	ErrType DBErrorType
	Message string
	Err     error
}

// Error returns the error message
func (e DBError) Error() string {
	message := e.Message
	if message == "" {
		if defaultMessage, ok := DBErrorDefaultMessages[e.ErrType]; ok {
			message = defaultMessage
		}
	}
	if e.Err != nil {
		message = strings.Join([]string{message, e.Err.Error()}, ": ")
	}
	return message
}

// Unwrap returns the wrapped error
func (e DBError) Unwrap() error {
	return e.Err
}
