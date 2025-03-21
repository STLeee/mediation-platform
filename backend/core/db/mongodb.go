package db

import (
	"context"
	"fmt"
	"time"

	"github.com/STLeee/mediation-platform/backend/core/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	DefaultMinPoolSize       = 1
	DefaultMaxPoolSize       = 10
	DefaultMaxConnIdleTime   = 10 * time.Second
	DefaultConnectionTimeout = 10 * time.Second
	DefaultTimeout           = 10 * time.Second
)

type MongoDBConfig struct {
	URI               string            `yaml:"uri"`
	TLS               bool              `yaml:"tls"`
	TLSConfig         *MongoDBTLSConfig `yaml:"tls_config"`
	MinPoolSize       int               `yaml:"min_pool_size"`
	MaxPoolSize       int               `yaml:"max_pool_size"`
	MaxConnIdleTime   time.Duration     `yaml:"max_conn_idle_time"`
	ConnectionTimeout time.Duration     `yaml:"connection_timeout"`
	Timeout           time.Duration     `yaml:"timeout"`
}

type MongoDBTLSConfig struct {
	CAFile   string `yaml:"ca_file"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

type MongoDB struct {
	client *mongo.Client
	cfg    *MongoDBConfig
}

func NewMongoDB(cfg *MongoDBConfig) (*MongoDB, error) {
	fmt.Printf("MongoDBConfig: %+v\n", cfg)

	// Set client options
	opt := options.Client().ApplyURI(cfg.URI)
	if cfg.TLS {
		if cfg.TLSConfig == nil {
			return nil, DBError{
				ErrType: DBErrorConfigError,
				Message: "TLS config is required",
			}
		}
		tlsConfig, err := utils.NewTLSConfig(cfg.TLSConfig.CAFile, cfg.TLSConfig.CertFile, cfg.TLSConfig.KeyFile)
		if err != nil {
			return nil, DBError{
				ErrType: DBErrorConfigError,
				Message: "failed to create TLS config",
				Err:     err,
			}
		}
		opt.SetTLSConfig(tlsConfig)
	}
	if cfg.MinPoolSize <= 0 {
		cfg.MinPoolSize = DefaultMinPoolSize
	}
	if cfg.MaxPoolSize <= 0 {
		cfg.MaxPoolSize = DefaultMaxPoolSize
	}
	if cfg.MaxConnIdleTime <= 0 {
		cfg.MaxConnIdleTime = DefaultMaxConnIdleTime
	}
	if cfg.ConnectionTimeout <= 0 {
		cfg.ConnectionTimeout = DefaultConnectionTimeout
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = DefaultTimeout
	}
	opt.SetMinPoolSize(uint64(cfg.MinPoolSize))
	opt.SetMaxPoolSize(uint64(cfg.MaxPoolSize))
	opt.SetMaxConnIdleTime(cfg.MaxConnIdleTime)
	opt.SetConnectTimeout(cfg.ConnectionTimeout)
	opt.SetTimeout(cfg.Timeout)

	// Connect to MongoDB
	client, err := mongo.Connect(opt)
	if err != nil {
		return nil, DBError{
			ErrType: DBErrorTypeServerError,
			Message: "failed to connect to MongoDB",
			Err:     err,
		}
	}

	mongoDB := &MongoDB{
		client: client,
		cfg:    cfg,
	}

	// Ping MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectionTimeout)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		mongoDB.Close()
		return nil, DBError{
			ErrType: DBErrorTypeServerError,
			Message: "failed to ping MongoDB",
			Err:     err,
		}
	}

	return mongoDB, nil
}

func (db *MongoDB) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), db.cfg.ConnectionTimeout)
	defer cancel()
	db.client.Disconnect(ctx)
}
