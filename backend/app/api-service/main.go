package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/STLeee/mediation-platform/backend/app/api-service/config"
	"github.com/STLeee/mediation-platform/backend/app/api-service/docs"
	"github.com/STLeee/mediation-platform/backend/app/api-service/middleware"
	"github.com/STLeee/mediation-platform/backend/app/api-service/router"
	coreAuth "github.com/STLeee/mediation-platform/backend/core/auth"
	coreService "github.com/STLeee/mediation-platform/backend/core/service"
)

func init() {
	// Load config
	cfg := loadConfig(config.DefaultConfigPath)

	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)
}

// @securityDefinitions.apikey TokenAuth
// @in header
// @name Authorization
func main() {
	// Get config
	cfg := config.GetConfig()

	// Init auth service
	authService := initAuthService(cfg)

	// Setup server
	engine := gin.Default()
	registerAPIRouters(engine, authService)

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

// Init auth service
func initAuthService(cfg *config.Config) coreAuth.BaseAuthService {
	authService, err := coreAuth.NewAuthService(context.Background(), &cfg.AuthService)
	if err != nil {
		panic(err)
	}

	return authService
}

// Register API routers
func registerAPIRouters(engine *gin.Engine, authService coreAuth.BaseAuthService) {
	// Register middleware
	engine.Use(middleware.CorsHandler())
	engine.Use(middleware.ErrorHandler())

	// Register API routers
	apiRouterGroup := engine.Group("/api")

	// Register health router
	healthRouterGroup := apiRouterGroup.Group("/health")
	router.RegisterHealthRouter(healthRouterGroup)

	// Register v1 router
	v1RouterGroup := apiRouterGroup.Group("/v1")

	// Register v1 user router
	userRouterGroup := v1RouterGroup.Group("/user")
	userRouterGroup.Use(middleware.TokenAuthHandler(authService))
	router.RegisterV1UserRouter(userRouterGroup)
}
