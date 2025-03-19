package main

import (
	"net/http"
	"os"
	"testing"

	"github.com/STLeee/mediation-platform/backend/core/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Run tests
	os.Exit(m.Run())
}

func TestApp(t *testing.T) {
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

func TestRegisterRouters(t *testing.T) {
	utils.TestEngineRouterRegister(t, func(engine *gin.Engine) {
		registerAPIRouters(engine, nil)
	}, []string{
		"/api/health/liveness",
		"/api/health/readiness",
		"/api/v1/user/:user_id",
	})
}
