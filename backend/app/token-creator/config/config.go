package config

import (
	"os"

	"gopkg.in/yaml.v3"

	coreAuth "github.com/STLeee/mediation-platform/backend/core/auth"
)

type Config struct {
	AuthService coreAuth.AuthServiceConfig `yaml:"auth_service"`
}

var cfg *Config

func LoadConfig(path string) (*Config, error) {
	loadedCfg := &Config{}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, loadedCfg)
	if err != nil {
		return nil, err
	}
	cfg = loadedCfg
	return cfg, nil
}

func GetConfig() *Config {
	return cfg
}
