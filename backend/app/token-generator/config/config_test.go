package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadAndGetConfig(t *testing.T) {
	// Create a temporary config file
	tempFile, err := os.CreateTemp("", "test_config_*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Write test data to the temporary config file
	configData := `
auth_service:
  firebase:
    project_id: mediation-platform-test
`
	_, err = tempFile.Write([]byte(configData))
	assert.NoError(t, err)
	tempFile.Close()

	// Load the config from the temporary file
	loadedCfg, err := LoadConfig(tempFile.Name())
	assert.NoError(t, err)
	assert.NotNil(t, loadedCfg)
	assert.Equal(t, "mediation-platform-test", loadedCfg.AuthService.FirebaseAuthConfig.ProjectID)

	// Ensure GetConfig returns the loaded config
	gotCfg := GetConfig()
	assert.Equal(t, loadedCfg, gotCfg)
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	// Load a non-existent config file
	_, err := LoadConfig("non_existent_file.yaml")
	assert.Error(t, err)
}
