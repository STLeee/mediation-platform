package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/STLeee/mediation-platform/backend/app/api-service/config"
	"github.com/STLeee/mediation-platform/backend/app/api-service/router"
)

func init() {
	// Load config
	cfg := loadConfig(config.DefaultConfigPath)

	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)
}

func main() {
	// Get config
	cfg := config.GetConfig()

	// Setup server
	engine := gin.Default()
	registerRouters(engine)

	// Run server
	engine.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}

// Load config
func loadConfig(path string) *config.Config {
	cfg, err := config.LoadConfig(path)
	if err != nil {
		panic(err)
	}

	return cfg
}

// Register routers
func registerRouters(engine *gin.Engine) {
	// Register health router
	healthRouter := engine.Group("/health")
	router.RegisterHealthRouter(healthRouter)
}
