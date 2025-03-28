package repository

import (
	"context"
	"testing"

	"github.com/STLeee/mediation-platform/backend/core/auth"
	"github.com/STLeee/mediation-platform/backend/core/model"
	"github.com/stretchr/testify/assert"
)

func TestSetDefaultUserCacheRepositoryConfig(t *testing.T) {
	testCases := []struct {
		name string
		cfg  *UserCacheRepositoryConfig
	}{
		{
			name: "default-config",
			cfg:  nil,
		},
		{
			name: "custom-config/no-keys",
			cfg:  &UserCacheRepositoryConfig{},
		},
		{
			name: "custom-config/empty-keys",
			cfg: &UserCacheRepositoryConfig{
				Keys: map[UserCacheRepositoryKeyName]*RedisCacheRepositoryKeyConfig{},
			},
		},
		{
			name: "custom-config/empty-key",
			cfg: &UserCacheRepositoryConfig{
				Keys: map[UserCacheRepositoryKeyName]*RedisCacheRepositoryKeyConfig{
					UserCacheRepositoryKeyNameAuthTokenUser: {},
				},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cfg := SetDefaultUserCacheRepositoryConfig(testCase.cfg)
			assert.NotNil(t, cfg)

			if testCase.cfg == nil {
				assert.Equal(t, DefaultUserCacheRepositoryConfig, cfg)
			} else {
				for _, key := range UserCacheRepositoryKeyNameList {
					if testCase.cfg.Keys[key] == nil {
						assert.Equal(t, DefaultUserCacheRepositoryConfig.Keys[key], cfg.Keys[key])
					} else {
						if testCase.cfg.Keys[key].KeyFormat == "" {
							assert.Equal(t, DefaultUserCacheRepositoryConfig.Keys[key].KeyFormat, cfg.Keys[key].KeyFormat)
						}
						if testCase.cfg.Keys[key].TTL == nil {
							assert.Equal(t, DefaultUserCacheRepositoryConfig.Keys[key].TTL, cfg.Keys[key].TTL)
						} else {
							assert.Equal(t, testCase.cfg.Keys[key].TTL.Expire, cfg.Keys[key].TTL.Expire)
						}
					}
				}
			}
		})
	}
}

func TestGenerateCacheKey(t *testing.T) {
	// Auth token user
	cacheKey := userRedisCacheRepository.generateAuthTokenUserCacheKey(auth.AuthServiceNameFirebase, "test-token")
	assert.Equal(t, "user:firebase:test-token", cacheKey)
}

func TestSetAuthTokenUser(t *testing.T) {
	ctx := context.Background()

	// Set user by auth token
	userRedisCacheRepository.SetAuthTokenUser(ctx, auth.AuthServiceNameFirebase, "test-token", localUsers[0])

	// Get user by auth token
	user, err := userRedisCacheRepository.GetAuthTokenUser(ctx, auth.AuthServiceNameFirebase, "test-token")
	assert.Nil(t, err)
	assertUser(t, localUsers[0], user)
}

func TestGetAuthTokenUser(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name         string
		authName     auth.AuthServiceName
		token        string
		exceptedUser *model.User
		exceptedErr  error
	}{
		{
			name:         "get-user",
			authName:     auth.AuthServiceNameFirebase,
			token:        "test-token",
			exceptedUser: localUsers[0],
			exceptedErr:  nil,
		},
		{
			name:         "not-found",
			authName:     auth.AuthServiceNameFirebase,
			token:        "not-found-token",
			exceptedUser: nil,
			exceptedErr: RepositoryError{
				ErrType: RepositoryErrorTypeRecordNotFound,
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.exceptedUser != nil {
				// Set user by auth token
				err := userRedisCacheRepository.SetAuthTokenUser(ctx, testCase.authName, testCase.token, testCase.exceptedUser)
				if err != nil {
					t.Fatal(err)
				}
			}

			// Get user by auth token
			user, err := userRedisCacheRepository.GetAuthTokenUser(ctx, testCase.authName, testCase.token)
			if testCase.exceptedUser != nil {
				assertUser(t, testCase.exceptedUser, user)
			}
			if testCase.exceptedErr != nil {
				assertError(t, testCase.exceptedErr, err)
			}
		})
	}
}
