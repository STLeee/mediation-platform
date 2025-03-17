package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	coreService "github.com/STLeee/mediation-platform/backend/core/service"
)

func TestLoadAndGetConfig(t *testing.T) {
	// Create a temporary config file
	tempFile, err := os.CreateTemp("", "test_config_*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Write test data to the temporary config file
	configData := `
server:
  port: 9090
  gin_mode: debug
service:
  name: test-service
  env: test
`
	_, err = tempFile.Write([]byte(configData))
	assert.NoError(t, err)
	tempFile.Close()

	// Load the config from the temporary file
	loadedCfg, err := LoadConfig(tempFile.Name())
	assert.NoError(t, err)
	assert.NotNil(t, loadedCfg)
	assert.Equal(t, 9090, loadedCfg.Server.Port)
	assert.Equal(t, "debug", loadedCfg.Server.GinMode)
	assert.Equal(t, "test-service", loadedCfg.Service.Name)
	assert.Equal(t, coreService.Testing, loadedCfg.Service.Environment)

	// Ensure GetConfig returns the loaded config
	gotCfg := GetConfig()
	assert.Equal(t, loadedCfg, gotCfg)
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	// Load a non-existent config file
	_, err := LoadConfig("non_existent_file.yaml")
	assert.Error(t, err)
}

func TestLoadConfig_InvalidConfig(t *testing.T) {
	// Create a temporary config file
	tempFile, err := os.CreateTemp("", "test_config_*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Write invalid test data to the temporary config file
	configData := `
server:
  port: invalid
service:
  name: test-service
  env: test
`
	_, err = tempFile.Write([]byte(configData))
	assert.NoError(t, err)
	tempFile.Close()

	// Load the config from the temporary file
	_, err = LoadConfig(tempFile.Name())
	assert.Error(t, err)
}
