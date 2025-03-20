package main

import (
	"net/http"
	"os"
	"testing"

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

	// Test health check
	response, err := http.Get("http://localhost:8080/api/health/liveness")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
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
	loadedCfg := loadConfig(tempFile.Name())
	assert.NotNil(t, loadedCfg)
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()

	// Load a non-existent config file
	loadConfig("non_existent_file.yaml")
}

func TestInitAuthService(t *testing.T) {
	// Test initAuthService function
	authService := initAuthService(&config.Config{
		AuthService: coreAuth.AuthServiceConfig{
			FirebaseAuthConfig: &coreAuth.FirebaseAuthConfig{
				ProjectID:    "test_project_id",
				KeyFile:      "test_key_file",
				EmulatorHost: "test_emulator_host",
			},
		},
	})
	assert.NotNil(t, authService)
}

func TestInitAuthService_NoConfig(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()

	initAuthService(&config.Config{})
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
