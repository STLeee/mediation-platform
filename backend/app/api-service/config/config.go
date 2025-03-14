package config

import (
	"os"

	"gopkg.in/yaml.v3"

	coreService "github.com/STLeee/mediation-platform/backend/core/service"
)

const DefaultConfigPath = "conf/app.conf.yaml"

type ServerConfig struct {
	Port int `yaml:"port" default:"8080"`
}

type ServiceConfig struct {
	Name        string                  `yaml:"name" default:"api-service"`
	Environment coreService.Environment `yaml:"env" default:"test"`
}

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Service ServiceConfig `yaml:"service"`
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
