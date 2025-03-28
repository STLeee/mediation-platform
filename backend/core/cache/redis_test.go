package cache

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var redisCache *RedisCache

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Init database
	var err error
	redisCache, err = NewRedisCache(ctx, LocalRedisCacheConfig)
	if err != nil {
		panic(err)
	}
	defer redisCache.Close()

	os.Exit(m.Run())
}

func TestNewRedisCacheFailed(t *testing.T) {
	ctx := context.Background()

	// Invalid address
	cfg := &RedisCacheConfig{
		Address:  "invalid-address",
		Username: "",
		Password: "",
		DB:       0,
	}
	_, err := NewRedisCache(ctx, cfg)
	assert.NotNil(t, err)
	assert.Equal(t, CacheErrorTypeServerError, err.(CacheError).ErrType)
	assert.Equal(t, "failed to ping Redis", err.(CacheError).Message)
	assert.NotNil(t, err.(CacheError).Err)
}

func TestRedisCache_Close(t *testing.T) {
	err := redisCache.Close()
	assert.Nil(t, err)
	err = redisCache.Close()
	assert.NotNil(t, err)
	assert.Equal(t, CacheErrorTypeServerError, err.(CacheError).ErrType)
	assert.Equal(t, "failed to close Redis connection", err.(CacheError).Message)
	assert.NotNil(t, err.(CacheError).Err)
}
