package repository

import (
	"context"
	"os"
	"testing"

	"github.com/STLeee/mediation-platform/backend/core/db"
	"github.com/STLeee/mediation-platform/backend/core/model"
	"github.com/STLeee/mediation-platform/backend/core/utils"
	"github.com/stretchr/testify/assert"
)

var userMongoDBRepository *UserMongoDBRepository

var localUsers = []*model.User{
	{
		UserID:      "67ac56653fee2207c557b99c",
		FirebaseUID: "W6WyRvhWhEarGHs7GV5unjVi8DYX",
		DisplayName: "TestingUser2",
		Email:       "testing2@mediation-platform.com",
		PhoneNumber: "",
		PhotoURL:    "",
		Disabled:    false,
		CreatedAt:   1742269900219,
		LastLoginAt: 1742269900219,
	},
}

func TestMain(m *testing.M) {
	// Connect to local MongoDB
	mongoDB, err := db.NewMongoDB(context.Background(), db.LocalMongoDBConfig)
	if err != nil {
		panic(err)
	}
	defer mongoDB.Close(context.Background())
	userMongoDBRepository = NewUserMongoDB(mongoDB, LocalMongoDBRepositoryConfigs[RepositoryNameUser])

	// Run tests
	os.Exit(m.Run())
}

func TestUserMongoDBRepository_GetUserByID(t *testing.T) {
	ctx := context.Background()

	// Test cases
	testCases := []struct {
		name         string
		userID       string
		expectedUser *model.User
		expectedErr  error
	}{
		{
			name:         "user-found",
			userID:       localUsers[0].UserID,
			expectedUser: localUsers[0],
			expectedErr:  nil,
		},
		{
			name:         "user-not-found",
			userID:       "aaaaaaaaaaaaaaaaaaaaaaaa",
			expectedUser: nil,
			expectedErr: RepositoryError{
				ErrType:    RepositoryErrorTypeRecordNotFound,
				Database:   LocalMongoDBRepositoryConfigs[RepositoryNameUser].Database,
				Collection: LocalMongoDBRepositoryConfigs[RepositoryNameUser].Collection,
				Message:    "record not found",
			},
		},
		{
			name:         "invalid-id",
			userID:       "invalid-id",
			expectedUser: nil,
			expectedErr: RepositoryError{
				ErrType:    RepositoryErrorTypeInvalidID,
				Database:   LocalMongoDBRepositoryConfigs[RepositoryNameUser].Database,
				Collection: LocalMongoDBRepositoryConfigs[RepositoryNameUser].Collection,
				Message:    "invalid ID",
			},
		},
	}

	// Run tests
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			user, err := userMongoDBRepository.GetUserByID(ctx, testCase.userID)
			if testCase.expectedUser != nil {
				assert.Nil(t, err)
				assert.Equal(t, utils.ToJSONString(testCase.expectedUser), utils.ToJSONString(user))
			} else {
				assert.Equal(t, testCase.expectedErr, err)
			}
		})
	}
}
