package auth

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"

	"github.com/STLeee/mediation-platform/backend/core/model"
)

var firebaseAuth *FirebaseAuth
var testUserList = []*model.UserInfo{
	{
		FirebaseUID:   "LRgwDJoRP7BCYJBNmNrNL4rxhvgR",
		DisplayName:   "TestingUser1",
		Email:         "testing1@mediation-platform.com",
		PhoneNumber:   "",
		PhotoURL:      "",
		Disabled:      false,
		EmailVerified: false,
	},
	{
		FirebaseUID:   "W6WyRvhWhEarGHs7GV5unjVi8DYX",
		DisplayName:   "TestingUser2",
		Email:         "testing2@mediation-platform.com",
		PhoneNumber:   "",
		PhotoURL:      "",
		Disabled:      false,
		EmailVerified: true,
	},
	{
		FirebaseUID:   "3fKQ3DyZhddm2H30J8ggTpsR35x2",
		DisplayName:   "TestingUser3",
		Email:         "testing3@mediation-platform.com",
		PhoneNumber:   "",
		PhotoURL:      "",
		Disabled:      true,
		EmailVerified: false,
	},
}

func TestMain(m *testing.M) {
	// Init Firebase
	var err error
	firebaseAuth, err = NewFirebaseAuth(context.Background(), &FirebaseAuthConfig{
		EmulatorHost: "localhost:9099",
		ProjectID:    "mediation-platform-test",
		KeyFile:      "",
	})
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

// createMockIDToken generates a mock Firebase ID token for testing purposes
func createMockIDToken(uid string) string {
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
		panic(err)
	}

	return tokenString
}

func TestFirebaseAuthenticateByToken(t *testing.T) {
	testCases := []struct {
		name    string
		uid     string
		token   string
		isValid bool
	}{
		{
			name:    "valid token",
			uid:     testUserList[0].FirebaseUID,
			token:   createMockIDToken(testUserList[0].FirebaseUID),
			isValid: true,
		},
		{
			name:    "invalid token",
			uid:     "testing",
			token:   "invalid-token",
			isValid: false,
		},
		{
			name:    "user not found",
			uid:     "not-found",
			token:   createMockIDToken("not-found"),
			isValid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			uid, err := firebaseAuth.AuthenticateByToken(context.Background(), testCase.token)
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

func TestGetUserInfo(t *testing.T) {
	testCases := []struct {
		name string
		uid  string
		want *model.UserInfo
	}{
		{
			name: testUserList[0].DisplayName,
			uid:  testUserList[0].FirebaseUID,
			want: testUserList[0],
		},
		{
			name: testUserList[1].DisplayName,
			uid:  testUserList[1].FirebaseUID,
			want: testUserList[1],
		},
		{
			name: testUserList[2].DisplayName,
			uid:  testUserList[2].FirebaseUID,
			want: testUserList[2],
		},
		{
			name: "invalid token",
			uid:  "invalid-uid",
			want: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			userInfo, err := firebaseAuth.GetUserInfo(context.Background(), testCase.uid)
			if testCase.want != nil {
				assert.Nil(t, err)
				assert.EqualValues(t, *testCase.want, *userInfo)
			} else {
				assert.NotNil(t, err)
				assert.Nil(t, userInfo)
			}
		})
	}
}
