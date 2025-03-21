package db

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBError(t *testing.T) {
	testCases := []struct {
		name     string
		errType  DBErrorType
		message  string
		err      error
		expected string
	}{
		{
			name:     "server-error/no-message",
			errType:  DBErrorTypeServerError,
			message:  "",
			err:      fmt.Errorf("test server error"),
			expected: "server error: test server error",
		},
		{
			name:     "server-error/with-message",
			errType:  DBErrorTypeServerError,
			message:  "test message",
			err:      fmt.Errorf("test server error"),
			expected: "test message: test server error",
		},
		{
			name:     "server-error/with-no-error",
			errType:  DBErrorTypeServerError,
			message:  "test message",
			err:      nil,
			expected: "test message",
		},
		{
			name:     "server-error/with-no-message-and-no-error",
			errType:  DBErrorTypeServerError,
			message:  "",
			err:      nil,
			expected: "server error",
		},
		{
			name:     "config-error",
			errType:  DBErrorConfigError,
			message:  "",
			err:      fmt.Errorf("test config error"),
			expected: "config error: test config error",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ae := DBError{
				ErrType: testCase.errType,
				Message: testCase.message,
				Err:     testCase.err,
			}
			assert.Equal(t, testCase.expected, ae.Error())
			assert.Equal(t, testCase.err, ae.Unwrap())
		})
	}
}
