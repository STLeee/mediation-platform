package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/STLeee/mediation-platform/backend/app/api-service/config"
	"github.com/STLeee/mediation-platform/backend/app/api-service/router"
)

func init() {
	// Load config
	_, err := config.LoadConfig(config.DefaultConfigPath)
	if err != nil {
		panic(err)
	}
}

func main() {
	// Get config
	cfg := config.GetConfig()

	// Register routers
	engine := gin.Default()
	registerRouters(engine)

	// Run server
	engine.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}

func registerRouters(engine *gin.Engine) {
	// Register health router
	healthRouter := engine.Group("/health")
	router.RegisterHealthRouter(healthRouter)
}
