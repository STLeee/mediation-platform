package auth

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/STLeee/mediation-platform/backend/core/model"
	"google.golang.org/api/option"
)

// FIREBASE_AUTH_NAME is the name of the firebase authentication
const FIREBASE_AUTH_NAME = "firebase"

// FirebaseAuthConfig struct for firebase authentication configuration
type FirebaseAuthConfig struct {
	ProjectID string `json:"project_id"`
	KeyFile   string `json:"key_file"`
}

// FirebaseAuthError struct for firebase authentication error
type FirebaseAuthError struct {
	AuthError
}

// Error returns the error message
func NewFirebaseAuthError(err error) FirebaseAuthError {
	return FirebaseAuthError{AuthError{FIREBASE_AUTH_NAME, err}}
}

// FirebaseAuth struct for firebase authentication
type FirebaseAuth struct {
	app        *firebase.App
	authClient *auth.Client
	cfg        *FirebaseAuthConfig
}

// NewFirebaseAuth creates a new FirebaseAuth struct
func NewFirebaseAuth(ctx context.Context, cfg *FirebaseAuthConfig) (*FirebaseAuth, error) {
	firebaseConfig := &firebase.Config{
		ProjectID: cfg.ProjectID,
	}
	options := []option.ClientOption{
		option.WithCredentialsFile(cfg.KeyFile),
	}
	app, err := firebase.NewApp(ctx, firebaseConfig, options...)
	if err != nil {
		return nil, NewFirebaseAuthError(err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, NewFirebaseAuthError(err)
	}

	return &FirebaseAuth{app, authClient, cfg}, nil
}

// AuthenticateByToken authenticates a user by token
func (firebaseAuth *FirebaseAuth) AuthenticateByToken(ctx context.Context, token string) (uid string, err error) {
	verifiedToken, err := firebaseAuth.authClient.VerifyIDToken(ctx, token)
	if err != nil {
		return "", NewFirebaseAuthError(err)
	}
	return verifiedToken.UID, nil
}

// GetUserInfo gets the user info by uid
func (firebaseAuth *FirebaseAuth) GetUserInfo(ctx context.Context, uid string) (*model.UserInfo, error) {
	userRecord, err := firebaseAuth.authClient.GetUser(ctx, uid)
	if err != nil {
		return nil, NewFirebaseAuthError(err)
	}
	return &model.UserInfo{
		UID:           userRecord.UID,
		DisplayName:   userRecord.DisplayName,
		Email:         userRecord.Email,
		PhoneNumber:   userRecord.PhoneNumber,
		PhotoURL:      userRecord.PhotoURL,
		Disabled:      userRecord.Disabled,
		EmailVerified: userRecord.EmailVerified,
	}, nil
}
