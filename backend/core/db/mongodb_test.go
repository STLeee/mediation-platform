package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAndCloseMongoDB(t *testing.T) {
	// Init database
	mongodb, err := NewMongoDB(context.Background(), LocalMongoDBConfig)
	assert.Nil(t, err)
	collection := mongodb.GetCollection("test-db", "test-collection")
	assert.NotNil(t, collection)

	defer mongodb.Close(context.Background())
}

func TestNewMongoDBError(t *testing.T) {
	testCases := []struct {
		name       string
		cfg        *MongoDBConfig
		errType    DBErrorType
		errMessage string
	}{
		{
			name: "tls-config-required",
			cfg: &MongoDBConfig{
				URI: "mongodb://localhost:27017",
				TLS: true,
			},
			errType:    DBErrorConfigError,
			errMessage: "TLS config is required",
		},
		{
			name: "tls-config-error",
			cfg: &MongoDBConfig{
				URI: "mongodb://localhost:27017",
				TLS: true,
				TLSConfig: &MongoDBTLSConfig{
					CAFile:   "invalid-ca.pem",
					CertFile: "invalid-client.pem",
					KeyFile:  "invalid-client.key",
				},
			},
			errType:    DBErrorConfigError,
			errMessage: "failed to create TLS config",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := NewMongoDB(context.Background(), testCase.cfg)
			assert.Equal(t, testCase.errType, err.(DBError).ErrType)
			assert.Equal(t, testCase.errMessage, err.(DBError).Message)
		})
	}
}
