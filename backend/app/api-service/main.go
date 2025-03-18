package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/STLeee/mediation-platform/backend/app/api-service/config"
	"github.com/STLeee/mediation-platform/backend/app/api-service/docs"
	"github.com/STLeee/mediation-platform/backend/app/api-service/middleware"
	"github.com/STLeee/mediation-platform/backend/app/api-service/router"
	coreService "github.com/STLeee/mediation-platform/backend/core/service"
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
	apiRouterGroup := engine.Group("/api")
	registerRouters(apiRouterGroup)

	// Swagger
	if cfg.Service.Environment == coreService.Testing {
		docs.SwaggerInfo.Title = "Mediation Platform API Service"
		docs.SwaggerInfo.Description = "API Service for Mediation Platform"
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		docs.SwaggerInfo.BasePath = "/api"
		docs.SwaggerInfo.Schemes = []string{"http", "https"}
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Run server
	engine.Run(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
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
func registerRouters(routerGroup *gin.RouterGroup) {
	// Middleware
	routerGroup.Use(middleware.Cors())

	// Register health router
	healthRouterGroup := routerGroup.Group("/health")
	router.RegisterHealthRouter(healthRouterGroup)
}
