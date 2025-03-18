package auth

import "fmt"

// AuthConfig struct for authentication configuration
type AuthConfig struct {
	FirebaseAuthConfig FirebaseAuthConfig `yaml:"firebase"`
}

// AuthError struct for authentication error
type AuthError struct {
	authName string
	err      error
}

// Error returns the error message
func (ae AuthError) Error() string {
	return fmt.Sprintf("[auth - %s] %s", ae.authName, ae.err.Error())
}

// Unwrap returns the wrapped error
func (ae AuthError) Unwrap() error {
	return ae.err
}
