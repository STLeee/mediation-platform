package cache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCacheError(t *testing.T) {
	testCases := []struct {
		name     string
		errType  CacheErrorType
		message  string
		err      error
		expected string
	}{
		{
			name:     "server-error/no-message",
			errType:  CacheErrorTypeServerError,
			message:  "",
			err:      fmt.Errorf("test server error"),
			expected: "server error: test server error",
		},
		{
			name:     "server-error/with-message",
			errType:  CacheErrorTypeServerError,
			message:  "test message",
			err:      fmt.Errorf("test server error"),
			expected: "test message: test server error",
		},
		{
			name:     "server-error/with-no-error",
			errType:  CacheErrorTypeServerError,
			message:  "test message",
			err:      nil,
			expected: "test message",
		},
		{
			name:     "server-error/with-no-message-and-no-error",
			errType:  CacheErrorTypeServerError,
			message:  "",
			err:      nil,
			expected: "server error",
		},
		{
			name:     "operation-error",
			errType:  CacheErrorTypeOperationError,
			message:  "",
			err:      fmt.Errorf("test error"),
			expected: "operation error: test error",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ae := CacheError{
				ErrType: testCase.errType,
				Message: testCase.message,
				Err:     testCase.err,
			}
			assert.Equal(t, testCase.expected, ae.Error())
			assert.Equal(t, testCase.err, ae.Unwrap())
		})
	}
}
