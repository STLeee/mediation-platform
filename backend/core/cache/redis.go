package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	DefaultPoolSize          = 10
	DefaultConnMaxIdleTime   = 30 * time.Minute
	DefaultConnectionTimeout = 5 * time.Second
	DefaultOperationTimeout  = 3 * time.Second
)

// RedisCacheConfig is the configuration for the RedisCache
type RedisCacheConfig struct {
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`

	PoolSize        int           `yaml:"pool_size"`
	ConnMaxIdleTime time.Duration `yaml:"idle_timeout"`
	DialTimeout     time.Duration `yaml:"dial_timeout"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
}

// LocalRedisCacheConfig is a config for local Redis
var LocalRedisCacheConfig = &RedisCacheConfig{
	Address:  "127.0.0.1:6379",
	Username: "",
	Password: "",
	DB:       0,
}

// RedisCache is a cache implementation using Redis
type RedisCache struct {
	redis.Client
	cfg *RedisCacheConfig
}

// NewRedisCache creates a new RedisCache instance
func NewRedisCache(ctx context.Context, cfg *RedisCacheConfig) (*RedisCache, error) {
	// Set options
	if cfg.PoolSize == 0 {
		cfg.PoolSize = DefaultPoolSize
	}
	if cfg.ConnMaxIdleTime == 0 {
		cfg.ConnMaxIdleTime = DefaultConnMaxIdleTime
	}
	if cfg.DialTimeout == 0 {
		cfg.DialTimeout = DefaultConnectionTimeout
	}
	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = DefaultOperationTimeout
	}
	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = DefaultOperationTimeout
	}
	opt := &redis.Options{
		Addr:            cfg.Address,
		Username:        cfg.Username,
		Password:        cfg.Password,
		DB:              cfg.DB,
		PoolSize:        cfg.PoolSize,
		ConnMaxIdleTime: cfg.ConnMaxIdleTime,
		DialTimeout:     cfg.DialTimeout,
		ReadTimeout:     cfg.ReadTimeout,
		WriteTimeout:    cfg.WriteTimeout,
	}

	// Create connection
	redisCache := &RedisCache{
		Client: *redis.NewClient(opt),
		cfg:    cfg,
	}

	// Ping
	_, err := redisCache.Ping(ctx).Result()
	if err != nil {
		return nil, CacheError{
			ErrType: CacheErrorTypeServerError,
			Message: "failed to ping Redis",
			Err:     err,
		}
	}

	return redisCache, nil
}

// Close closes the Redis connection
func (redisCache *RedisCache) Close() error {
	err := redisCache.Client.Close()
	if err != nil {
		return CacheError{
			ErrType: CacheErrorTypeServerError,
			Message: "failed to close Redis connection",
			Err:     err,
		}
	}
	return nil
}
