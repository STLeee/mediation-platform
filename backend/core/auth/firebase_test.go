package auth

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

var firebaseAuth *FirebaseAuth

func TestMain(m *testing.M) {
	// Set Firebase Auth emulator host environment variable
	emulatorHost := "localhost:9099"
	if err := os.Setenv("FIREBASE_AUTH_EMULATOR_HOST", emulatorHost); err != nil {
		panic(err)
	}

	// Init Firebase
	var err error
	firebaseAuth, err = NewFirebaseAuth(context.Background(), &FirebaseAuthConfig{
		ProjectID: "mediation-platform-test",
		KeyFile:   "",
	})
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

// createMockIDToken generates a mock Firebase ID token for testing purposes
func createMockIDToken(uid string) (string, error) {
	// Define the token claims
	now := time.Now().Unix()
	exp := now + 60
	claims := jwt.MapClaims{
		"iss": fmt.Sprintf("https://securetoken.google.com/%s", firebaseAuth.cfg.ProjectID),
		"aud": firebaseAuth.cfg.ProjectID,
		"uid": uid,
		"sub": uid,
		"iat": now,
		"exp": exp,
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a dummy secret
	secret := []byte("testing-secret")
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func TestFirebaseAuthenticateByToken(t *testing.T) {
	testCases := []struct {
		name    string
		uid     string
		isValid bool
	}{
		{
			name:    "valid token",
			uid:     "LRgwDJoRP7BCYJBNmNrNL4rxhvgR",
			isValid: true,
		},
		{
			name:    "invalid token",
			uid:     "invalid-uid",
			isValid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockToken, err := createMockIDToken(testCase.uid)
			assert.Nil(t, err)
			assert.NotEmpty(t, mockToken)

			uid, err := firebaseAuth.AuthenticateByToken(context.Background(), mockToken)
			if testCase.isValid {
				assert.Nil(t, err)
				assert.Equal(t, testCase.uid, uid)
			} else {
				assert.NotNil(t, err)
				assert.Empty(t, uid)
			}
		})
	}
}
