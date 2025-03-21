package repository

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepositoryError(t *testing.T) {
	testCases := []struct {
		name       string
		errType    RepositoryErrorType
		database   string
		collection string
		message    string
		err        error
		expected   string
	}{
		{
			name:       "server-error/no-message",
			errType:    RepositoryErrorTypeServerError,
			database:   "test-db",
			collection: "test-collection",
			message:    "",
			err:        fmt.Errorf("test error"),
			expected:   "test-db/test-collection: server error: test error",
		},
		{
			name:       "server-error/with-message",
			errType:    RepositoryErrorTypeServerError,
			database:   "test-db",
			collection: "test-collection",
			message:    "test message",
			err:        fmt.Errorf("test error"),
			expected:   "test-db/test-collection: test message: test error",
		},
		{
			name:       "server-error/with-no-error",
			errType:    RepositoryErrorTypeServerError,
			database:   "test-db",
			collection: "test-collection",
			message:    "test message",
			err:        nil,
			expected:   "test-db/test-collection: test message",
		},
		{
			name:       "server-error/with-no-message-and-no-error",
			errType:    RepositoryErrorTypeServerError,
			database:   "test-db",
			collection: "test-collection",
			message:    "",
			err:        nil,
			expected:   "test-db/test-collection: server error",
		},
		{
			name:       "record-not-found/no-message",
			errType:    RepositoryErrorTypeRecordNotFound,
			database:   "test-db",
			collection: "test-collection",
			message:    "",
			err:        fmt.Errorf("test error"),
			expected:   "test-db/test-collection: record not found: test error",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ae := RepositoryError{
				ErrType:    testCase.errType,
				Database:   testCase.database,
				Collection: testCase.collection,
				Message:    testCase.message,
				Err:        testCase.err,
			}
			assert.Equal(t, testCase.expected, ae.Error())
			assert.Equal(t, testCase.err, ae.Unwrap())
		})
	}
}
