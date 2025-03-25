package main

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/STLeee/mediation-platform/backend/app/api-service/config"
	coreAuth "github.com/STLeee/mediation-platform/backend/core/auth"
	"github.com/STLeee/mediation-platform/backend/core/utils"
)

func TestApp(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Test main function
	go main()

	// Test health check (retry 10 times with 1 second interval)
	for i := 0; i < 10; i++ {
		response, err := http.Get("http://localhost:8080/api/health/liveness")
		if err == nil && response.StatusCode == http.StatusOK {
			return
		}
		time.Sleep(1 * time.Second)
	}
	assert.Fail(t, "Failed to start the server")
}

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	tempFile, err := os.CreateTemp("", "test_config_*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Write test data to the temporary config file
	_, err = tempFile.Write([]byte(""))
	assert.NoError(t, err)
	tempFile.Close()

	// Load the config from the temporary file
	loadedCfg, err := loadConfig(tempFile.Name())
	assert.NoError(t, err)
	assert.NotNil(t, loadedCfg)
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	// Load a non-existent config file
	loadedCfg, err := loadConfig("non_existent_file.yaml")
	assert.Error(t, err)
	assert.Nil(t, loadedCfg)
}

func TestInitAuthService(t *testing.T) {
	testCases := []struct {
		name    string
		config  *config.Config
		isError bool
	}{
		{
			name: "valid-config",
			config: &config.Config{
				AuthService: coreAuth.AuthServiceConfig{
					FirebaseAuthConfig: &coreAuth.FirebaseAuthConfig{
						ProjectID:    "test_project_id",
						KeyFile:      "test_key_file",
						EmulatorHost: "test_emulator_host",
					},
				},
			},
			isError: false,
		},
		{
			name:    "no-config",
			config:  &config.Config{},
			isError: true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			authService, err := initAuthService(testCase.config)
			if !testCase.isError {
				assert.NoError(t, err)
				assert.NotNil(t, authService)
			} else {
				assert.Error(t, err)
			}
		})
	}

	// Test initAuthService function
	authService, err := initAuthService(&config.Config{
		AuthService: coreAuth.AuthServiceConfig{
			FirebaseAuthConfig: &coreAuth.FirebaseAuthConfig{
				ProjectID:    "test_project_id",
				KeyFile:      "test_key_file",
				EmulatorHost: "test_emulator_host",
			},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, authService)
}

func TestInitMongoDB(t *testing.T) {
	testCases := []struct {
		name    string
		config  *config.Config
		isError bool
	}{
		{
			name:    "no-config",
			config:  &config.Config{},
			isError: true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mongoDB, err := initMongoDB(testCase.config)
			if !testCase.isError {
				assert.NoError(t, err)
				assert.NotNil(t, mongoDB)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestRegisterRouters(t *testing.T) {
	utils.TestEngineRouterRegister(t, func(engine *gin.Engine) {
		registerAPIRouters(engine, nil)
	}, []string{
		"/api/health/liveness",
		"/api/health/readiness",
		"/api/v1/user/:user_id",
	})
}
