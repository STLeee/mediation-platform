package auth

import (
	"context"
	"strings"

	"github.com/STLeee/mediation-platform/backend/core/model"
)

// AuthServiceConfig struct for authentication service configuration
type AuthServiceConfig struct {
	FirebaseAuthConfig *FirebaseAuthConfig `yaml:"firebase"`
}

type AuthServiceErrorType string

const (
	AuthServiceErrorTypeServerError  AuthServiceErrorType = "server_error"
	AuthServiceErrorTypeTokenInvalid AuthServiceErrorType = "token_invalid"
	AuthServiceErrorTypeUserNotFound AuthServiceErrorType = "user_not_found"
)

var AuthServiceErrorDefaultMessages = map[AuthServiceErrorType]string{
	AuthServiceErrorTypeServerError:  "server error",
	AuthServiceErrorTypeTokenInvalid: "token is invalid",
	AuthServiceErrorTypeUserNotFound: "user not found",
}

// AuthServiceError struct for authentication service error
type AuthServiceError struct {
	ErrType AuthServiceErrorType
	Message string
	Err     error
}

// Error returns the error message
func (e AuthServiceError) Error() string {
	message := e.Message
	if message == "" {
		if defaultMessage, ok := AuthServiceErrorDefaultMessages[e.ErrType]; ok {
			message = defaultMessage
		}
	}
	if e.Err != nil {
		message = strings.Join([]string{message, e.Err.Error()}, ": ")
	}
	return message
}

// Unwrap returns the wrapped error
func (e AuthServiceError) Unwrap() error {
	return e.Err
}

// BaseAuthService interface for authentication service
type BaseAuthService interface {
	AuthenticateByToken(ctx context.Context, token string) (uid string, err error)
	GetUserInfo(ctx context.Context, uid string) (*model.User, error)
}

// NewAuthService creates a new authentication service
func NewAuthService(ctx context.Context, cfg *AuthServiceConfig) (BaseAuthService, error) {
	if cfg != nil {
		if cfg.FirebaseAuthConfig != nil {
			return NewFirebaseAuth(ctx, cfg.FirebaseAuthConfig)
		}
	}
	return nil, AuthServiceError{
		ErrType: AuthServiceErrorTypeServerError,
		Message: "no authentication service is configured",
	}
}
