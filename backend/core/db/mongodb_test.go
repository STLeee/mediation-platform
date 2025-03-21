package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testMongoDBConfig = &MongoDBConfig{
	URI: "mongodb://admin:pass@127.0.0.1:27017/?directConnection=true",
	TLS: true,
	TLSConfig: &MongoDBTLSConfig{
		CAFile:   "../../mongodb/tls/test-ca.pem",
		CertFile: "../../mongodb/tls/test-client.pem",
		KeyFile:  "../../mongodb/tls/mongodb-test-client.key",
	},
}

func TestMongoDB(t *testing.T) {
	// Init database
	mongodb, err := NewMongoDB(testMongoDBConfig)
	assert.Nil(t, err)
	defer mongodb.Close()
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
			_, err := NewMongoDB(testCase.cfg)
			assert.Equal(t, testCase.errType, err.(DBError).ErrType)
			assert.Equal(t, testCase.errMessage, err.(DBError).Message)
		})
	}
}
