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
		errType  AuthServiceErrorType
		message  string
		err      error
		expected string
	}{
		{
			errType:  AuthServiceErrorTypeServerError,
			message:  "",
			err:      fmt.Errorf("test server error"),
			expected: "server error: test server error",
		},
		{
			errType:  AuthServiceErrorTypeServerError,
			message:  "test message",
			err:      fmt.Errorf("test server error"),
			expected: "test message: test server error",
		},
		{
			errType:  AuthServiceErrorTypeServerError,
			message:  "test message",
			err:      nil,
			expected: "test message",
		},
		{
			errType:  AuthServiceErrorTypeServerError,
			message:  "",
			err:      nil,
			expected: "server error",
		},
		{
			errType:  AuthServiceErrorTypeTokenInvalid,
			message:  "",
			err:      fmt.Errorf("test token invalid"),
			expected: "token is invalid: test token invalid",
		},
		{
			errType:  AuthServiceErrorTypeUserNotFound,
			message:  "",
			err:      fmt.Errorf("test user not found"),
			expected: "user not found: test user not found",
		},
	}

	for _, testCase := range testCases {
		t.Run(string(testCase.errType), func(t *testing.T) {
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
		id       string
		cfg      *AuthServiceConfig
		expected any
	}{
		{
			id: "firebase",
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
			id:  "no config",
			cfg: nil,
			expected: AuthServiceError{
				ErrType: AuthServiceErrorTypeServerError,
				Message: "no authentication service is configured",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			authService, err := NewAuthService(context.Background(), testCase.cfg)
			if err == nil {
				excepted := reflect.TypeOf(testCase.expected).Elem().Name()
				actual := reflect.TypeOf(authService).Elem().Name()
				assert.Equal(t, excepted, actual)
			} else {
				t.Log(err)
				assert.Equal(t, testCase.expected, err)
			}
		})
	}
}
