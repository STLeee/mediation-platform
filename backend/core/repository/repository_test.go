package repository

import (
	"context"
	"fmt"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/STLeee/mediation-platform/backend/core/cache"
	"github.com/STLeee/mediation-platform/backend/core/db"
	"github.com/STLeee/mediation-platform/backend/core/model"
	"github.com/STLeee/mediation-platform/backend/core/utils"
	"github.com/stretchr/testify/assert"
)

var (
	userMongoDBRepository    *UserMongoDBRepository
	userRedisCacheRepository *UserRedisCacheRepository
)

var localUsers = []*model.User{
	{
		UserID:      "000000000000000000000001",
		FirebaseUID: "LRgwDJoRP7BCYJBNmNrNL4rxhvgR",
		DisplayName: "TestingUser1",
		Email:       "testing1@mediation-platform.com",
		PhoneNumber: "",
		PhotoURL:    "",
		Disabled:    false,
		CreatedAt:   time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
		LastLoginAt: time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      "000000000000000000000002",
		FirebaseUID: "W6WyRvhWhEarGHs7GV5unjVi8DYX",
		DisplayName: "TestingUser2",
		Email:       "testing2@mediation-platform.com",
		PhoneNumber: "",
		PhotoURL:    "",
		Disabled:    false,
		CreatedAt:   time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
		LastLoginAt: time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      "000000000000000000000003",
		FirebaseUID: "3fKQ3DyZhddm2H30J8ggTpsR35x2",
		DisplayName: "TestingUser3",
		Email:       "testing3@mediation-platform.com",
		PhoneNumber: "",
		PhotoURL:    "",
		Disabled:    true,
		CreatedAt:   time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
		LastLoginAt: time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
	},
}

func assertError(t *testing.T, expectedErr, actualErr error) {
	assert.ErrorAs(t, actualErr, &expectedErr)
	if _, ok := actualErr.(RepositoryError); ok {
		assert.Equal(t, expectedErr.(RepositoryError).ErrType, actualErr.(RepositoryError).ErrType)
		assert.Equal(t, expectedErr.(RepositoryError).Database, actualErr.(RepositoryError).Database)
		assert.Equal(t, expectedErr.(RepositoryError).Collection, actualErr.(RepositoryError).Collection)
	}
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

func TestMain(m *testing.M) {
	// Connect to local MongoDB
	mongoDB, err := db.NewMongoDB(context.Background(), db.LocalMongoDBConfig)
	if err != nil {
		panic(err)
	}
	defer mongoDB.Close()

	// Connect to local Redis
	redis, err := cache.NewRedisCache(context.Background(), cache.LocalRedisCacheConfig)
	if err != nil {
		panic(err)
	}
	defer redis.Close()

	userMongoDBRepository = NewUserMongoDBRepository(mongoDB, LocalRepositoryConfigs.UserDB)
	userRedisCacheRepository = NewUserRedisCacheRepository(redis, nil)

	// Run tests
	os.Exit(m.Run())
}

func TestRepositoryError(t *testing.T) {
	testCases := []struct {
		name       string
		errType    RepositoryErrorType
		database   string
		collection string
		message    string
		err        error
		expected   string
	}{
		{
			name:       "server-error/no-message",
			errType:    RepositoryErrorTypeServerError,
			database:   "test-db",
			collection: "test-collection",
			message:    "",
			err:        fmt.Errorf("test error"),
			expected:   "test-db/test-collection: server error: test error",
		},
		{
			name:       "server-error/with-message",
			errType:    RepositoryErrorTypeServerError,
			database:   "test-db",
			collection: "test-collection",
			message:    "test message",
			err:        fmt.Errorf("test error"),
			expected:   "test-db/test-collection: test message: test error",
		},
		{
			name:       "server-error/with-no-error",
			errType:    RepositoryErrorTypeServerError,
			database:   "test-db",
			collection: "test-collection",
			message:    "test message",
			err:        nil,
			expected:   "test-db/test-collection: test message",
		},
		{
			name:       "server-error/with-no-message-and-no-error",
			errType:    RepositoryErrorTypeServerError,
			database:   "test-db",
			collection: "test-collection",
			message:    "",
			err:        nil,
			expected:   "test-db/test-collection: server error",
		},
		{
			name:       "record-not-found/no-message",
			errType:    RepositoryErrorTypeRecordNotFound,
			database:   "test-db",
			collection: "test-collection",
			message:    "",
			err:        fmt.Errorf("test error"),
			expected:   "test-db/test-collection: record not found: test error",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ae := RepositoryError{
				ErrType:    testCase.errType,
				Database:   testCase.database,
				Collection: testCase.collection,
				Message:    testCase.message,
				Err:        testCase.err,
			}
			assert.Equal(t, testCase.expected, ae.Error())
			assert.Equal(t, testCase.err, ae.Unwrap())
		})
	}
}
