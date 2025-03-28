package repository

import (
	"context"
	"time"

	"github.com/STLeee/mediation-platform/backend/core/auth"
	"github.com/STLeee/mediation-platform/backend/core/cache"
	"github.com/STLeee/mediation-platform/backend/core/model"
	"github.com/redis/go-redis/v9"
)

// UserCacheRepository is an interface for user cache repository
type UserCacheRepository interface {
	SetAuthTokenUser(ctx context.Context, authName auth.AuthServiceName, token string, user *model.User) error
	GetAuthTokenUser(ctx context.Context, authName auth.AuthServiceName, token string) (*model.User, error)
}

// UserCacheKeyPrefix is a prefix for user cache key
const UserCacheKeyPrefix = "user"

// UserCacheRepositoryKeyName is a key name for user cache repository
type UserCacheRepositoryKeyName string

const (
	UserCacheRepositoryKeyNameAuthTokenUser UserCacheRepositoryKeyName = "auth_token_user"
)

var UserCacheRepositoryKeyNameList = []UserCacheRepositoryKeyName{
	UserCacheRepositoryKeyNameAuthTokenUser,
}

// UserCacheRepositoryConfig is a configuration for UserCacheRepository
type UserCacheRepositoryConfig struct {
	Keys map[UserCacheRepositoryKeyName]*RedisCacheRepositoryKeyConfig
}

// DefaultUserCacheRepositoryConfig is a default configuration for UserCacheRepository
var DefaultUserCacheRepositoryConfig = &UserCacheRepositoryConfig{
	Keys: map[UserCacheRepositoryKeyName]*RedisCacheRepositoryKeyConfig{
		UserCacheRepositoryKeyNameAuthTokenUser: {
			KeyFormat: "{auth_name}:{token}",
			TTL: &RedisCacheRepositoryKeyTTLConfig{
				Expire:          1 * time.Hour,
				MaxRandomOffset: 5 * time.Minute,
			},
		},
	},
}

// UserRedisCacheRepository
type UserRedisCacheRepository struct {
	RedisCacheRepository
	cfg *UserCacheRepositoryConfig
}

func SetDefaultUserCacheRepositoryConfig(cfg *UserCacheRepositoryConfig) *UserCacheRepositoryConfig {
	if cfg == nil {
		return DefaultUserCacheRepositoryConfig
	} else {
		if cfg.Keys == nil {
			cfg.Keys = DefaultUserCacheRepositoryConfig.Keys
		} else {
			for _, key := range UserCacheRepositoryKeyNameList {
				if cfg.Keys[key] == nil {
					cfg.Keys[key] = DefaultUserCacheRepositoryConfig.Keys[key]
				} else {
					if cfg.Keys[key].KeyFormat == "" {
						cfg.Keys[key].KeyFormat = DefaultUserCacheRepositoryConfig.Keys[key].KeyFormat
					}
					if cfg.Keys[key].TTL == nil {
						cfg.Keys[key].TTL = DefaultUserCacheRepositoryConfig.Keys[key].TTL
					}
				}
			}
		}
	}
	return cfg
}

// NewUserRedisCacheRepository creates a new UserRedisCacheRepository
func NewUserRedisCacheRepository(redisCache *cache.RedisCache, cfg *UserCacheRepositoryConfig) *UserRedisCacheRepository {
	cfg = SetDefaultUserCacheRepositoryConfig(cfg)
	return &UserRedisCacheRepository{
		RedisCacheRepository: *NewRedisCacheRepository(redisCache),
		cfg:                  cfg,
	}
}

// generateAuthTokenUserCacheKey generates cache key for user by auth token
func (repo *UserRedisCacheRepository) generateAuthTokenUserCacheKey(authName auth.AuthServiceName, token string) string {
	cacheKeyCfg := repo.cfg.Keys[UserCacheRepositoryKeyNameAuthTokenUser]
	return UserCacheKeyPrefix + ":" + cacheKeyCfg.GenerateCacheKey(map[string]string{
		"{auth_name}": string(authName),
		"{token}":     token,
	})
}

// SetAuthTokenUser sets user by auth token
func (repo *UserRedisCacheRepository) SetAuthTokenUser(ctx context.Context, authName auth.AuthServiceName, token string, user *model.User) error {
	cacheKeyCfg := repo.cfg.Keys[UserCacheRepositoryKeyNameAuthTokenUser]
	cacheKey := repo.generateAuthTokenUserCacheKey(authName, token)
	cacheValue, err := repo.ConvertToJSON(user)
	if err != nil {
		return err
	}
	ttl := cacheKeyCfg.TTL.GenerateTTL()

	err = repo.Set(ctx, cacheKey, cacheValue, ttl).Err()
	if err != nil {
		return RepositoryError{
			ErrType: RepositoryErrorTypeServerError,
			Message: "failed to set user by auth token",
			Err:     err,
		}
	}
	return nil
}

// GetAuthTokenUser gets user by auth token
func (repo *UserRedisCacheRepository) GetAuthTokenUser(ctx context.Context, authName auth.AuthServiceName, token string) (*model.User, error) {
	cacheKey := repo.generateAuthTokenUserCacheKey(authName, token)
	cacheValue, err := repo.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, RepositoryError{
				ErrType: RepositoryErrorTypeRecordNotFound,
				Message: "user not found by auth token",
			}
		}
		return nil, RepositoryError{
			ErrType: RepositoryErrorTypeServerError,
			Message: "failed to get user by auth token",
			Err:     err,
		}
	}
	var user model.User
	err = repo.RevertFromJSON(cacheValue, &user)
	if err != nil {
		return nil, RepositoryError{
			ErrType: RepositoryErrorTypeServerError,
			Message: "failed to convert user from JSON",
			Err:     err,
		}
	}
	return &user, nil
}
