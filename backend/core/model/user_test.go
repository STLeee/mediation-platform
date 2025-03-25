package model

import (
	"testing"

	"github.com/STLeee/mediation-platform/backend/core/utils"
	"github.com/stretchr/testify/assert"
)

func TestUserInMongoDB(t *testing.T) {
	testCases := []struct {
		name    string
		user    *User
		isValid bool
	}{
		{
			name: "valid-user",
			user: &User{
				UserID:      "5f4b8f1f9d1e4b0001f3f3b1",
				DisplayName: "display-name",
				Email:       "email",
				PhoneNumber: "phone-number",
			},
			isValid: true,
		},
		{
			name: "empty-user-id",
			user: &User{
				UserID:      "",
				DisplayName: "display-name",
				Email:       "email",
				PhoneNumber: "phone-number",
			},
			isValid: true,
		},
		{
			name: "invalid-user-id",
			user: &User{
				UserID:      "invalid-id",
				DisplayName: "display-name",
				Email:       "email",
				PhoneNumber: "phone-number",
			},
			isValid: false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			userInMongoDB, err := NewUserInMongoDB(testCase.user)
			if testCase.isValid {
				if err != nil {
					t.Fatal(err)
				}

				if testCase.user.UserID != "" {
					exceptedUserID := testCase.user.UserID

					// Check if the user ID is converted to ObjectID
					assert.Equal(t, utils.ConvertStringToObjectID(exceptedUserID), userInMongoDB.ID)

					// Check if the ObjectID is set to the user
					userInMongoDB.User.UserID = ""
					userInMongoDB.SetupDataFromDocument()
					assert.Equal(t, exceptedUserID, userInMongoDB.User.UserID)
				} else {
					// Check if the ObjectID is set to the user
					userInMongoDB.SetupDataFromDocument()
					assert.NotEmpty(t, userInMongoDB.User.UserID)
				}
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}
