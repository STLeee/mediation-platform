package db

import (
	"context"
	"time"

	"github.com/STLeee/mediation-platform/backend/core/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Default values
const (
	DefaultMinPoolSize       = 1
	DefaultMaxPoolSize       = 10
	DefaultMaxConnIdleTime   = 10 * time.Second
	DefaultConnectionTimeout = 10 * time.Second
	DefaultTimeout           = 10 * time.Second
)

// LocalMongoDBConfig is a config for local MongoDB
var LocalMongoDBConfig = &MongoDBConfig{
	URI: "mongodb://admin:pass@127.0.0.1:27017/?directConnection=true",
	TLS: true,
	TLSConfig: &MongoDBTLSConfig{
		CAFile:   "../../mongodb/tls/test-ca.pem",
		CertFile: "../../mongodb/tls/test-client.pem",
		KeyFile:  "../../mongodb/tls/mongodb-test-client.key",
	},
}

// MongoDBConfig is a config for MongoDB
type MongoDBConfig struct {
	URI string `yaml:"uri"`

	TLS       bool              `yaml:"tls"`
	TLSConfig *MongoDBTLSConfig `yaml:"tls_config"`

	MinPoolSize int `yaml:"min_pool_size"`
	MaxPoolSize int `yaml:"max_pool_size"`

	MaxConnIdleTime   time.Duration `yaml:"max_conn_idle_time"`
	ConnectionTimeout time.Duration `yaml:"connection_timeout"`
	Timeout           time.Duration `yaml:"timeout"`
}

// MongoDBTLSConfig is a config for MongoDB TLS
type MongoDBTLSConfig struct {
	CAFile   string `yaml:"ca_file"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

// MongoDB is a MongoDB client
type MongoDB struct {
	mongo.Client
	cfg *MongoDBConfig
}

// NewMongoDB creates a new MongoDB
func NewMongoDB(ctx context.Context, cfg *MongoDBConfig) (*MongoDB, error) {
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
		Client: *client,
		cfg:    cfg,
	}

	// Ping MongoDB
	timeoutCtx, cancel := context.WithTimeout(ctx, cfg.ConnectionTimeout)
	defer cancel()
	if err := client.Ping(timeoutCtx, nil); err != nil {
		mongoDB.Close()
		return nil, DBError{
			ErrType: DBErrorTypeServerError,
			Message: "failed to ping MongoDB",
			Err:     err,
		}
	}

	return mongoDB, nil
}

// Close closes the MongoDB client
func (db *MongoDB) Close() {
	ctx := context.Background()
	timeoutCtx, cancel := context.WithTimeout(ctx, db.cfg.ConnectionTimeout)
	defer cancel()
	db.Disconnect(timeoutCtx)
}
