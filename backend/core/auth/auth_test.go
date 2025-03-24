package auth

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthError(t *testing.T) {
	testCases := []struct {
		name     string
		errType  AuthServiceErrorType
		message  string
		err      error
		expected string
	}{
		{
			name:     "server-error/no-message",
			errType:  AuthServiceErrorTypeServerError,
			message:  "",
			err:      fmt.Errorf("test server error"),
			expected: "server error: test server error",
		},
		{
			name:     "server-error/with-message",
			errType:  AuthServiceErrorTypeServerError,
			message:  "test message",
			err:      fmt.Errorf("test server error"),
			expected: "test message: test server error",
		},
		{
			name:     "server-error/with-no-error",
			errType:  AuthServiceErrorTypeServerError,
			message:  "test message",
			err:      nil,
			expected: "test message",
		},
		{
			name:     "server-error/with-no-message-and-no-error",
			errType:  AuthServiceErrorTypeServerError,
			message:  "",
			err:      nil,
			expected: "server error",
		},
		{
			name:     "token-invalid",
			errType:  AuthServiceErrorTypeTokenInvalid,
			message:  "",
			err:      fmt.Errorf("test token invalid"),
			expected: "token is invalid: test token invalid",
		},
		{
			name:     "user-not-found",
			errType:  AuthServiceErrorTypeUserNotFound,
			message:  "",
			err:      fmt.Errorf("test user not found"),
			expected: "user not found: test user not found",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ae := AuthServiceError{
				ErrType: testCase.errType,
				Message: testCase.message,
				Err:     testCase.err,
			}
			assert.Equal(t, testCase.expected, ae.Error())
			assert.Equal(t, testCase.err, ae.Unwrap())
		})
	}
}

func TestNewAuthService(t *testing.T) {
	testCases := []struct {
		name     string
		cfg      *AuthServiceConfig
		expected any
	}{
		{
			name: "firebase",
			cfg: &AuthServiceConfig{
				FirebaseAuthConfig: &FirebaseAuthConfig{
					ProjectID:    "test_project_id",
					KeyFile:      "test_key_file",
					EmulatorHost: "test_emulator_host",
				},
			},
			expected: &FirebaseAuth{},
		},
		{
			name: "no-config",
			cfg:  nil,
			expected: AuthServiceError{
				ErrType: AuthServiceErrorTypeServerError,
				Message: "no authentication service is configured",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			authService, err := NewAuthService(context.Background(), testCase.cfg)
			if err == nil {
				expected := reflect.TypeOf(testCase.expected).Elem().Name()
				actual := reflect.TypeOf(authService).Elem().Name()
				assert.Equal(t, expected, actual)
			} else {
				t.Log(err)
				assert.Equal(t, testCase.expected, err)
			}
		})
	}
}
