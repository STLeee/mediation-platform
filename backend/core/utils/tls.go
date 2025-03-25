package utils

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
)

// NewTLSConfig creates a new TLS config
func NewTLSConfig(caFile, certFile, keyFile string) (*tls.Config, error) {
	// Read CA certificate
	caCert, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		return nil, errors.New("failed to append ca certificate to pool")
	}

	// Load client certificate
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
	}, nil
}
