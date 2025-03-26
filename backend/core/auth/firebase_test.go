package auth

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/STLeee/mediation-platform/backend/core/model"
	"github.com/STLeee/mediation-platform/backend/core/utils"
)

var firebaseAuth *FirebaseAuth
var localUsers = []*model.User{
	{
		FirebaseUID: "LRgwDJoRP7BCYJBNmNrNL4rxhvgR",
		DisplayName: "TestingUser1",
		Email:       "testing1@mediation-platform.com",
		PhoneNumber: "",
		PhotoURL:    "",
		Disabled:    false,
	},
	{
		FirebaseUID: "W6WyRvhWhEarGHs7GV5unjVi8DYX",
		DisplayName: "TestingUser2",
		Email:       "testing2@mediation-platform.com",
		PhoneNumber: "",
		PhotoURL:    "",
		Disabled:    false,
	},
	{
		FirebaseUID: "3fKQ3DyZhddm2H30J8ggTpsR35x2",
		DisplayName: "TestingUser3",
		Email:       "testing3@mediation-platform.com",
		PhoneNumber: "",
		PhotoURL:    "",
		Disabled:    true,
	},
}

func TestMain(m *testing.M) {
	// Connect to local Firebase
	var err error
	firebaseAuth, err = NewFirebaseAuth(context.Background(), LocalFirebaseAuthConfig)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestFirebase_GetName(t *testing.T) {
	assert.Equal(t, AuthServiceNameFirebase, firebaseAuth.GetName())
}

func TestFirebase_AuthenticateByToken(t *testing.T) {
	testCases := []struct {
		name        string
		uid         string
		token       string
		expectedErr error
	}{
		{
			name:        "valid-token",
			uid:         localUsers[0].FirebaseUID,
			token:       utils.CreateMockFirebaseIDToken(firebaseAuth.cfg.ProjectID, localUsers[0].FirebaseUID),
			expectedErr: nil,
		},
		{
			name:        "invalid-token",
			uid:         "testing",
			token:       "invalid-token",
			expectedErr: AuthServiceError{ErrType: AuthServiceErrorTypeTokenInvalid},
		},
		{
			name:        "user-not-found",
			uid:         "not-found",
			token:       utils.CreateMockFirebaseIDToken(firebaseAuth.cfg.ProjectID, "not-found"),
			expectedErr: AuthServiceError{ErrType: AuthServiceErrorTypeUserNotFound},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			uid, err := firebaseAuth.AuthenticateByToken(context.Background(), testCase.token)
			if testCase.expectedErr == nil {
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, testCase.uid, uid)
			} else {
				assert.ErrorAs(t, err, &testCase.expectedErr)
				assert.Empty(t, uid)
				if _, ok := err.(AuthServiceError); ok {
					assert.Equal(t, testCase.expectedErr.(AuthServiceError).ErrType, err.(AuthServiceError).ErrType)
				}
			}
		})
	}
}

func TestFirebase_GetUserInfo(t *testing.T) {
	testCases := []struct {
		name string
		uid  string
		want *model.User
		err  error
	}{
		{
			name: localUsers[0].DisplayName,
			uid:  localUsers[0].FirebaseUID,
			want: localUsers[0],
		},
		{
			name: localUsers[1].DisplayName,
			uid:  localUsers[1].FirebaseUID,
			want: localUsers[1],
		},
		{
			name: localUsers[2].DisplayName,
			uid:  localUsers[2].FirebaseUID,
			want: localUsers[2],
		},
		{
			name: "user-not-found",
			uid:  "invalid-uid",
			err: AuthServiceError{
				ErrType: AuthServiceErrorTypeUserNotFound,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			userInfo, err := firebaseAuth.GetUserInfo(context.Background(), testCase.uid)
			if testCase.want != nil {
				assert.Equal(t, *testCase.want, *userInfo)
			} else {
				assert.Nil(t, userInfo)
			}
			assert.Equal(t, testCase.err, err)
		})
	}
}
