package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpStatusCodeError(t *testing.T) {
	testCases := []struct {
		id         string
		statusCode int
		message    string
		err        error
		expected   string
	}{
		{
			id:         "400",
			statusCode: 400,
			message:    "",
			err:        fmt.Errorf("test-error"),
			expected:   "Bad Request",
		},
		{
			id:         "400 with message",
			statusCode: 400,
			message:    "test message",
			err:        fmt.Errorf("test-error"),
			expected:   "test message",
		},
		{
			id:         "400 with message and no error",
			statusCode: 400,
			message:    "test message",
			err:        nil,
			expected:   "test message",
		},
		{
			id:         "400 with no error",
			statusCode: 400,
			message:    "",
			err:        nil,
			expected:   "Bad Request",
		},
		{
			id:         "401",
			statusCode: 401,
			message:    "",
			err:        fmt.Errorf("test-error"),
			expected:   "Unauthorized",
		},
		{
			id:         "404",
			statusCode: 404,
			message:    "",
			err:        fmt.Errorf("test-error"),
			expected:   "Not Found",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
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
