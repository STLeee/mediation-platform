package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/STLeee/mediation-platform/backend/app/token-generator/config"
	"github.com/STLeee/mediation-platform/backend/core/utils"
)

func main() {
	// Disable timestamp
	log.SetFlags(0)

	// Parse arguments
	uid := flag.String("uid", "", "Firebase UID")
	configPath := flag.String("config", "../api-service/conf/app.conf.yaml", "Config file path")
	flag.Parse()

	// Check arguments
	if *uid == "" {
		log.Fatalf("uid is required")
	}

	// Load config
	cfg, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create token
	token, err := generateToken(*uid, cfg)
	if err != nil {
		log.Fatalf("Failed to create token: %v", err)
	}

	// Print token
	log.Print(token)
}

// Load config
func loadConfig(path string) (*config.Config, error) {
	cfg, err := config.LoadConfig(path)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// Create token
func generateToken(uid string, cfg *config.Config) (string, error) {
	if cfg.AuthService.FirebaseAuthConfig != nil {
		return utils.GenerateMockFirebaseIDToken(cfg.AuthService.FirebaseAuthConfig.ProjectID, uid), nil
	}
	return "", fmt.Errorf("auth config is not set or not supported")
}
