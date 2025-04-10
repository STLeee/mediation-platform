package auth

import (
	"context"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	firebaseErrorutils "firebase.google.com/go/v4/errorutils"
	"github.com/STLeee/mediation-platform/backend/core/model"
	"google.golang.org/api/option"
)

var LocalFirebaseAuthConfig = &FirebaseAuthConfig{
	EmulatorHost: "localhost:9099",
	ProjectID:    "mediation-platform-test",
	KeyFile:      "",
}

// FirebaseAuthConfig struct for firebase authentication configuration
type FirebaseAuthConfig struct {
	ProjectID    string `yaml:"project_id"`
	KeyFile      string `yaml:"key_file"`
	EmulatorHost string `yaml:"emulator_host"`
}

// FirebaseAuth struct for firebase authentication
type FirebaseAuth struct {
	app        *firebase.App
	authClient *auth.Client
	cfg        *FirebaseAuthConfig
}

// NewFirebaseAuth creates a new FirebaseAuth struct
func NewFirebaseAuth(ctx context.Context, cfg *FirebaseAuthConfig) (*FirebaseAuth, error) {
	// Set Firebase Auth emulator host environment variable
	if err := os.Setenv("FIREBASE_AUTH_EMULATOR_HOST", cfg.EmulatorHost); err != nil {
		return nil, AuthServiceError{
			ErrType: AuthServiceErrorTypeServerError,
			Message: "failed to set firebase auth emulator host",
			Err:     err,
		}
	}

	// Create Firebase app
	firebaseConfig := &firebase.Config{
		ProjectID: cfg.ProjectID,
	}
	options := []option.ClientOption{
		option.WithCredentialsFile(cfg.KeyFile),
	}
	app, err := firebase.NewApp(ctx, firebaseConfig, options...)
	if err != nil {
		return nil, AuthServiceError{
			ErrType: AuthServiceErrorTypeServerError,
			Message: "failed to create firebase app",
			Err:     err,
		}
	}

	// Create Firebase Auth client
	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, AuthServiceError{
			ErrType: AuthServiceErrorTypeServerError,
			Message: "failed to create firebase auth client",
			Err:     err,
		}
	}

	return &FirebaseAuth{app, authClient, cfg}, nil
}

// GetName returns the authentication service name
func (firebaseAuth *FirebaseAuth) GetName() AuthServiceName {
	return AuthServiceNameFirebase
}

// AuthenticateByToken authenticates a user by token
func (firebaseAuth *FirebaseAuth) AuthenticateByToken(ctx context.Context, token string) (uid string, err error) {
	verifiedToken, err := firebaseAuth.authClient.VerifyIDToken(ctx, token)
	if err != nil {
		if firebaseErrorutils.IsInvalidArgument(err) {
			return "", AuthServiceError{ErrType: AuthServiceErrorTypeTokenInvalid}
		} else if firebaseErrorutils.IsNotFound(err) {
			return "", AuthServiceError{ErrType: AuthServiceErrorTypeUserNotFound}
		}
		return "", AuthServiceError{
			ErrType: AuthServiceErrorTypeServerError,
			Message: "failed to verify token",
			Err:     err,
		}
	}
	return verifiedToken.UID, nil
}

// GetUserInfo gets user info
func (firebaseAuth *FirebaseAuth) GetUserInfo(ctx context.Context, uid string) (user *model.User, err error) {
	userRecord, err := firebaseAuth.authClient.GetUser(ctx, uid)
	if err != nil {
		if firebaseErrorutils.IsNotFound(err) {
			return nil, AuthServiceError{ErrType: AuthServiceErrorTypeUserNotFound}
		}
		return nil, AuthServiceError{
			ErrType: AuthServiceErrorTypeServerError,
			Message: "failed to get user info",
			Err:     err,
		}
	}
	user = &model.User{
		FirebaseUID: userRecord.UID,
		DisplayName: userRecord.DisplayName,
		Email:       userRecord.Email,
		PhoneNumber: userRecord.PhoneNumber,
		PhotoURL:    userRecord.PhotoURL,
		Disabled:    userRecord.Disabled,
	}
	return user, nil
}
