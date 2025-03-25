package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/STLeee/mediation-platform/backend/app/token-creator/config"
	coreAuth "github.com/STLeee/mediation-platform/backend/core/auth"
	"github.com/STLeee/mediation-platform/backend/core/utils"
)

func TestApp(t *testing.T) {
	testCases := []struct {
		name       string
		uid        string
		configPath string
		configData string
		isSuccess  bool
	}{
		{
			name:      "valid-uid-and-default-config",
			uid:       "000000000000000000000001",
			isSuccess: true,
		},
		{
			name:      "empty-uid",
			uid:       "",
			isSuccess: false,
		},
		{
			name:       "config-not-found",
			uid:        "000000000000000000000001",
			configPath: "temp",
			configData: "",
			isSuccess:  false,
		},
		{
			name:       "auth-config-not-set",
			uid:        "000000000000000000000001",
			configPath: "temp",
			configData: `...`,
			isSuccess:  false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Set arguments
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()
			newArgs := []string{"cmd"}
			if testCase.uid != "" {
				newArgs = append(newArgs, fmt.Sprintf("-uid=%s", testCase.uid))
			}
			if testCase.configPath != "" {
				configPath := testCase.configPath
				if testCase.configPath == "temp" {
					// Create a temporary config file
					tempFile, err := os.CreateTemp("", "test_config_*.yaml")
					assert.NoError(t, err)
					defer os.Remove(tempFile.Name())
					configPath = tempFile.Name()

					// Write test data to the temporary config file
					_, err = tempFile.Write([]byte(testCase.configData))
					assert.NoError(t, err)
					tempFile.Close()
				}

				newArgs = append(newArgs, fmt.Sprintf("-config=%s", configPath))
			}
			os.Args = newArgs

			// Capture stdout
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			// Not success case
			if !testCase.isSuccess {
				assert.Panics(t, main)
				return
			}

			// Success case
			main()
			token := buf.String()
			t.Log(token)
			assert.NotEmpty(t, token)

			// Decode token
			decodedToken, err := utils.DecodeMockFirebaseIDToken(token)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, testCase.uid, decodedToken["uid"])
			t.Log(decodedToken)
		})
	}
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

func TestCreateToken(t *testing.T) {
	testCases := []struct {
		name    string
		uid     string
		config  *config.Config
		isError bool
	}{
		{
			name: "valid-config",
			uid:  "000000000000000000000001",
			config: &config.Config{
				AuthService: coreAuth.AuthServiceConfig{
					FirebaseAuthConfig: &coreAuth.FirebaseAuthConfig{
						ProjectID: "test_project_id",
					},
				},
			},
			isError: false,
		},
		{
			name:    "no-config",
			uid:     "000000000000000000000001",
			config:  &config.Config{},
			isError: true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			token, err := createToken(testCase.uid, testCase.config)
			if !testCase.isError {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
