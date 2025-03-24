package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// CreateMockFirebaseIDToken creates a mock Firebase ID token for testing
func CreateMockFirebaseIDToken(projectID, uid string) string {
	// Define the token claims
	now := time.Now().Unix()
	exp := now + 60
	claims := jwt.MapClaims{
		"iss": fmt.Sprintf("https://securetoken.google.com/%s", projectID),
		"aud": projectID,
		"uid": uid,
		"sub": uid,
		"iat": now,
		"exp": exp,
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a dummy secret
	secret := []byte("testing-secret")
	tokenString, _ := token.SignedString(secret)

	return tokenString
}
