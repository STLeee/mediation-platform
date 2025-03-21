package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTLSConfig(t *testing.T) {
	testCases := []struct {
		name      string
		caFile    string
		certFile  string
		keyFile   string
		isSuccess bool
	}{
		{
			name:      "success",
			caFile:    "../../mongodb/tls/test-ca.pem",
			certFile:  "../../mongodb/tls/test-client.pem",
			keyFile:   "../../mongodb/tls/mongodb-test-client.key",
			isSuccess: true,
		},
		{
			name:      "ca-file-not-exist",
			caFile:    "invalid-ca.pem",
			certFile:  "../../mongodb/tls/test-client.pem",
			keyFile:   "../../mongodb/tls/mongodb-test-client.key",
			isSuccess: false,
		},
		{
			name:      "ca-file-invalid",
			caFile:    "../../mongodb/tls/mongodb-test-client.key",
			certFile:  "../../mongodb/tls/test-client.pem",
			keyFile:   "../../mongodb/tls/mongodb-test-client.key",
			isSuccess: false,
		},
		{
			name:      "cert-file-not-exist",
			caFile:    "../../mongodb/tls/test-ca.pem",
			certFile:  "invalid-client.pem",
			keyFile:   "../../mongodb/tls/mongodb-test-client.key",
			isSuccess: false,
		},
		{
			name:      "key-file-not-exist",
			caFile:    "../../mongodb/tls/test-ca.pem",
			certFile:  "../../mongodb/tls/test-client.pem",
			keyFile:   "invalid-client.key",
			isSuccess: false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tlsConfig, err := NewTLSConfig(testCase.caFile, testCase.certFile, testCase.keyFile)
			if testCase.isSuccess {
				assert.Nil(t, err)
				assert.NotNil(t, tlsConfig)
				assert.NotNil(t, tlsConfig.RootCAs)
				assert.NotNil(t, tlsConfig.Certificates)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}
