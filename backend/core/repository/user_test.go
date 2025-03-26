package repository

import (
	"context"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/STLeee/mediation-platform/backend/core/auth"
	"github.com/STLeee/mediation-platform/backend/core/db"
	"github.com/STLeee/mediation-platform/backend/core/model"
	"github.com/STLeee/mediation-platform/backend/core/utils"
	"github.com/stretchr/testify/assert"
)

var userMongoDBRepository *UserMongoDBRepository

var localUsers = []*model.User{
	{
		UserID:      "000000000000000000000001",
		FirebaseUID: "LRgwDJoRP7BCYJBNmNrNL4rxhvgR",
		DisplayName: "TestingUser1",
		Email:       "testing1@mediation-platform.com",
		PhoneNumber: "",
		PhotoURL:    "",
		Disabled:    false,
	},
	{
		UserID:      "000000000000000000000002",
		FirebaseUID: "W6WyRvhWhEarGHs7GV5unjVi8DYX",
		DisplayName: "TestingUser2",
		Email:       "testing2@mediation-platform.com",
		PhoneNumber: "",
		PhotoURL:    "",
		Disabled:    false,
	},
	{
		UserID:      "000000000000000000000003",
		FirebaseUID: "3fKQ3DyZhddm2H30J8ggTpsR35x2",
		DisplayName: "TestingUser3",
		Email:       "testing3@mediation-platform.com",
		PhoneNumber: "",
		PhotoURL:    "",
		Disabled:    true,
	},
	{
		// Not in db
		UserID:      "",
		FirebaseUID: "kEnwA5bzGJrkEAnO2atgv6Fbbc2X",
		DisplayName: "TestingUser4",
		Email:       "testing4@mediation-platform.com",
		PhoneNumber: "",
		PhotoURL:    "",
		Disabled:    false,
	},
}

func assertUserWithMap(t *testing.T, expectedMap map[string]any, actualUser *model.User) {
	timestampFields := []string{"CreatedAt", "UpdatedAt", "LastLoginAt"}
	actualMap, err := utils.ConvertStructToMap(actualUser)
	if err != nil {
		t.Fatal(err)
	}
	for field, value := range expectedMap {
		if slices.Contains(timestampFields, field) {
			assert.True(t, utils.SimplyValidTimestamp(actualMap[field].(time.Time)))
		} else {
			assert.Equal(t, value, actualMap[field])
		}
	}
}

func assertUser(t *testing.T, expectedUser, actualUser *model.User) {
	expectedMap, err := utils.ConvertStructToMap(expectedUser)
	if err != nil {
		t.Fatal(err)
	}
	assertUserWithMap(t, expectedMap, actualUser)
}

func assertError(t *testing.T, expectedErr, actualErr error) {
	assert.ErrorAs(t, actualErr, &expectedErr)
	if _, ok := actualErr.(RepositoryError); ok {
		assert.Equal(t, expectedErr.(RepositoryError).ErrType, actualErr.(RepositoryError).ErrType)
		assert.Equal(t, expectedErr.(RepositoryError).Database, actualErr.(RepositoryError).Database)
		assert.Equal(t, expectedErr.(RepositoryError).Collection, actualErr.(RepositoryError).Collection)
	}
}

func TestMain(m *testing.M) {
	// Connect to local Firebase
	var err error
	firebaseAuth, err := auth.NewFirebaseAuth(context.Background(), auth.LocalFirebaseAuthConfig)
	if err != nil {
		panic(err)
	}

	// Connect to local MongoDB
	mongoDB, err := db.NewMongoDB(context.Background(), db.LocalMongoDBConfig)
	if err != nil {
		panic(err)
	}
	defer mongoDB.Close()

	userMongoDBRepository = NewUserMongoDBRepository(firebaseAuth, mongoDB, LocalMongoDBRepositoryConfigs[RepositoryNameUser])

	// Run tests
	os.Exit(m.Run())
}

func TestGetUserByToken(t *testing.T) {
	ctx := context.Background()

	// Test cases
	testCases := []struct {
		name         string
		token        string
		expectedUser *model.User
		expectedErr  error
		isNotInDB    bool
	}{
		{
			name:         "user-found",
			token:        utils.CreateMockFirebaseIDToken(auth.LocalFirebaseAuthConfig.ProjectID, localUsers[0].FirebaseUID),
			expectedUser: localUsers[0],
			expectedErr:  nil,
		},
		{
			name:        "invalid-token",
			token:       "invalid-token",
			expectedErr: auth.AuthServiceError{ErrType: auth.AuthServiceErrorTypeTokenInvalid},
		},
		{
			name:        "user-not-found-from-auth",
			token:       utils.CreateMockFirebaseIDToken(auth.LocalFirebaseAuthConfig.ProjectID, "not-found-auth-uid"),
			expectedErr: auth.AuthServiceError{ErrType: auth.AuthServiceErrorTypeUserNotFound},
		},
		{
			name:         "user-not-in-db",
			token:        utils.CreateMockFirebaseIDToken(auth.LocalFirebaseAuthConfig.ProjectID, localUsers[3].FirebaseUID),
			expectedUser: localUsers[3],
			isNotInDB:    true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			user, err := userMongoDBRepository.GetUserByToken(ctx, testCase.token)
			if testCase.expectedUser != nil {
				if err != nil {
					t.Fatal(err)
				}

				// Defer clean up
				if testCase.isNotInDB {
					assert.Empty(t, testCase.expectedUser.UserID)
					testCase.expectedUser.UserID = user.UserID
					defer func() {
						err := userMongoDBRepository.DeleteUserByID(ctx, user.UserID)
						if err != nil {
							t.Fatal(err)
						}
					}()
				}

				assertUser(t, testCase.expectedUser, user)
			} else {
				assert.ErrorAs(t, err, &testCase.expectedErr)
				if _, ok := err.(RepositoryError); ok {
					assert.Equal(t, testCase.expectedErr.(auth.AuthServiceError).ErrType, err.(RepositoryError).ErrType)
				}
			}
		})
	}
}

func TestUserMongoDBRepository_CreateUser(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name        string
		user        *model.User
		expectedErr error
	}{
		{
			name: "insert-user/with_user_id",
			user: &model.User{
				UserID:      "111111111111111111111111",
				DisplayName: "test-create-user",
				Email:       "test-create-user@mediation-platform.com",
			},
			expectedErr: nil,
		},
		{
			name: "insert-user/without_user_id",
			user: &model.User{
				DisplayName: "test-create-user",
				Email:       "test-create-user@mediation-platform.com",
			},
			expectedErr: nil,
		},
		{
			name: "insert-user/invalid-user_id",
			user: &model.User{
				UserID:      "invalid-id",
				DisplayName: "test-create-user",
			},
			expectedErr: RepositoryError{
				ErrType:    RepositoryErrorTypeInvalidID,
				Database:   LocalMongoDBRepositoryConfigs[RepositoryNameUser].Database,
				Collection: LocalMongoDBRepositoryConfigs[RepositoryNameUser].Collection,
			},
		},
		{
			name: "insert-user-duplicate",
			user: localUsers[0],
			expectedErr: RepositoryError{
				ErrType:    RepositoryErrorTypeServerError,
				Database:   LocalMongoDBRepositoryConfigs[RepositoryNameUser].Database,
				Collection: LocalMongoDBRepositoryConfigs[RepositoryNameUser].Collection,
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			userID, err := userMongoDBRepository.CreateUser(ctx, testCase.user)
			if testCase.expectedErr == nil {
				if err != nil {
					t.Fatal(err)
				}

				// Defer clean up
				defer func() {
					err := userMongoDBRepository.DeleteUserByID(ctx, userID)
					if err != nil {
						t.Fatal(err)
					}
				}()

				// Check user id
				if testCase.user.UserID != "" {
					assert.Equal(t, testCase.user.UserID, userID)
				} else {
					assert.NotEmpty(t, userID)
					testCase.user.UserID = userID
				}

				// Check if the user is created
				user, err := userMongoDBRepository.GetUserByID(ctx, userID)
				assert.Nil(t, err)
				if err != nil {
					assertUser(t, testCase.user, user)
				}
			} else {
				assert.ErrorAs(t, err, &testCase.expectedErr)
				assert.Empty(t, userID)
				if _, ok := err.(RepositoryError); ok {
					assert.Equal(t, testCase.expectedErr.(RepositoryError).ErrType, err.(RepositoryError).ErrType)
					assert.Equal(t, testCase.expectedErr.(RepositoryError).Database, err.(RepositoryError).Database)
					assert.Equal(t, testCase.expectedErr.(RepositoryError).Collection, err.(RepositoryError).Collection)
				}
			}
		})
	}
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
			},
		},
	}

	// Run tests
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			user, err := userMongoDBRepository.GetUserByID(ctx, testCase.userID)
			if testCase.expectedUser != nil {
				if err != nil {
					t.Fatal(err)
				}
				assertUser(t, testCase.expectedUser, user)
			} else {
				assertError(t, testCase.expectedErr, err)
			}
		})
	}
}

func TestUserMongoDBRepository_UpdateUser(t *testing.T) {
	ctx := context.Background()

	// Insert test user
	testUser := &model.User{
		UserID:      "111111111111111111111111",
		DisplayName: "test-update-user",
		Email:       "test-update-user@mediation-platform.com",
	}
	userID, err := userMongoDBRepository.CreateUser(ctx, testUser)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := userMongoDBRepository.DeleteUserByID(ctx, userID)
		if err != nil {
			t.Fatal(err)
		}
	}()

	testCases := []struct {
		name        string
		userID      string
		updateData  map[string]any
		expectedErr error
	}{
		{
			name:   "update-user",
			userID: userID,
			updateData: map[string]any{
				"display_name": "test-update-user-updated",
			},
			expectedErr: nil,
		},
		{
			name:   "update-user-not-found",
			userID: "aaaaaaaaaaaaaaaaaaaaaaaa",
			updateData: map[string]any{
				"display_name": "test-update-user-updated",
			},
			expectedErr: RepositoryError{
				ErrType:    RepositoryErrorTypeRecordNotFound,
				Database:   LocalMongoDBRepositoryConfigs[RepositoryNameUser].Database,
				Collection: LocalMongoDBRepositoryConfigs[RepositoryNameUser].Collection,
			},
		},
		{
			name:   "invalid-id",
			userID: "invalid-id",
			updateData: map[string]any{
				"display_name": "test-update-user-updated",
			},
			expectedErr: RepositoryError{
				ErrType:    RepositoryErrorTypeInvalidID,
				Database:   LocalMongoDBRepositoryConfigs[RepositoryNameUser].Database,
				Collection: LocalMongoDBRepositoryConfigs[RepositoryNameUser].Collection,
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := userMongoDBRepository.UpdateUserByID(ctx, testCase.userID, testCase.updateData)
			if testCase.expectedErr == nil {
				if err != nil {
					t.Fatal(err)
				}

				// Check if the user is updated
				user, err := userMongoDBRepository.GetUserByID(ctx, testCase.userID)
				assert.Nil(t, err)
				if err != nil {
					assertUserWithMap(t, testCase.updateData, user)
				}
			} else {
				assert.ErrorAs(t, err, &testCase.expectedErr)
				if _, ok := err.(RepositoryError); ok {
					assert.Equal(t, testCase.expectedErr.(RepositoryError).ErrType, err.(RepositoryError).ErrType)
					assert.Equal(t, testCase.expectedErr.(RepositoryError).Database, err.(RepositoryError).Database)
					assert.Equal(t, testCase.expectedErr.(RepositoryError).Collection, err.(RepositoryError).Collection)
				}
			}
		})
	}
}

func TestUserMongoDBRepository_DeleteUser(t *testing.T) {
	ctx := context.Background()

	testUser := &model.User{
		UserID:      "111111111111111111111111",
		DisplayName: "test-update-user",
		Email:       "test-update-user@mediation-platform.com",
	}

	testCases := []struct {
		name        string
		userID      string
		expectedErr error
	}{
		{
			name:        "delete-user",
			userID:      testUser.UserID,
			expectedErr: nil,
		},
		{
			name:   "delete-user-not-found",
			userID: "aaaaaaaaaaaaaaaaaaaaaaaa",
			expectedErr: RepositoryError{
				ErrType:    RepositoryErrorTypeRecordNotFound,
				Database:   LocalMongoDBRepositoryConfigs[RepositoryNameUser].Database,
				Collection: LocalMongoDBRepositoryConfigs[RepositoryNameUser].Collection,
			},
		},
		{
			name:   "invalid-id",
			userID: "invalid-id",
			expectedErr: RepositoryError{
				ErrType:    RepositoryErrorTypeInvalidID,
				Database:   LocalMongoDBRepositoryConfigs[RepositoryNameUser].Database,
				Collection: LocalMongoDBRepositoryConfigs[RepositoryNameUser].Collection,
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Insert test user
			userID, err := userMongoDBRepository.CreateUser(ctx, testUser)
			if err != nil {
				t.Fatal(err)
			}

			err = userMongoDBRepository.DeleteByID(ctx, testCase.userID)
			if testCase.expectedErr == nil {
				if err != nil {
					t.Fatal(err)
				}

				// Check if the user is deleted
				_, err := userMongoDBRepository.GetUserByID(ctx, testCase.userID)
				assert.NotNil(t, err)
			} else {
				// Defer clean up
				defer func() {
					err := userMongoDBRepository.DeleteUserByID(ctx, userID)
					if err != nil {
						t.Fatal(err)
					}
				}()

				assert.ErrorAs(t, err, &testCase.expectedErr)
				if _, ok := err.(RepositoryError); ok {
					assert.Equal(t, testCase.expectedErr.(RepositoryError).ErrType, err.(RepositoryError).ErrType)
					assert.Equal(t, testCase.expectedErr.(RepositoryError).Database, err.(RepositoryError).Database)
					assert.Equal(t, testCase.expectedErr.(RepositoryError).Collection, err.(RepositoryError).Collection)
				}
			}
		})
	}
}
