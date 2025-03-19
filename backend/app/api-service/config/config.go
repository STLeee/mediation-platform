package config

import (
	"os"

	"gopkg.in/yaml.v3"

	coreAuth "github.com/STLeee/mediation-platform/backend/core/auth"
	coreService "github.com/STLeee/mediation-platform/backend/core/service"
)

const DefaultConfigPath = "conf/app.conf.yaml"

type ServerConfig struct {
	Host    string `yaml:"host" default:"localhost"`
	Port    int    `yaml:"port" default:"8080"`
	GinMode string `yaml:"gin_mode" default:"release" validate:"oneof=debug release test"`
}

type ServiceConfig struct {
	Name        string                         `yaml:"name" default:"api-service"`
	Environment coreService.ServiceEnvironment `yaml:"env" default:"test" validate:"oneof=test stag prod"`
}

type Config struct {
	Server      ServerConfig               `yaml:"server"`
	Service     ServiceConfig              `yaml:"service"`
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
