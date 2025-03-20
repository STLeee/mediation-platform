package model

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpStatusCodeError(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		message    string
		err        error
		expected   string
	}{
		{
			name:       "bad-request",
			statusCode: http.StatusBadRequest,
			message:    "",
			err:        fmt.Errorf("test-error"),
			expected:   http.StatusText(http.StatusBadRequest),
		},
		{
			name:       "bad-request/with-message",
			statusCode: http.StatusBadRequest,
			message:    "test message",
			err:        fmt.Errorf("test-error"),
			expected:   "test message",
		},
		{
			name:       "bad-request/with-no-error",
			statusCode: http.StatusBadRequest,
			message:    "test message",
			err:        nil,
			expected:   "test message",
		},
		{
			name:       "bad-request/with-no-message-and-no-error",
			statusCode: http.StatusBadRequest,
			message:    "",
			err:        nil,
			expected:   http.StatusText(http.StatusBadRequest),
		},
		{
			name:       "unauthorized",
			statusCode: http.StatusUnauthorized,
			message:    "",
			err:        fmt.Errorf("test-error"),
			expected:   http.StatusText(http.StatusUnauthorized),
		},
		{
			name:       "not-found",
			statusCode: http.StatusNotFound,
			message:    "",
			err:        fmt.Errorf("test-error"),
			expected:   http.StatusText(http.StatusNotFound),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ae := HttpStatusCodeError{
				StatusCode: testCase.statusCode,
				Message:    testCase.message,
				Err:        testCase.err,
			}
			assert.Equal(t, testCase.expected, ae.Error())
			assert.Equal(t, testCase.err, ae.Unwrap())
		})
	}
}
