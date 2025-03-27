package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/STLeee/mediation-platform/backend/app/api-service/config"
	"github.com/STLeee/mediation-platform/backend/app/api-service/docs"
	"github.com/STLeee/mediation-platform/backend/app/api-service/middleware"
	"github.com/STLeee/mediation-platform/backend/app/api-service/router"
	coreAuth "github.com/STLeee/mediation-platform/backend/core/auth"
	coreCache "github.com/STLeee/mediation-platform/backend/core/cache"
	coreDB "github.com/STLeee/mediation-platform/backend/core/db"
	coreRepository "github.com/STLeee/mediation-platform/backend/core/repository"
	coreService "github.com/STLeee/mediation-platform/backend/core/service"
)

// @securityDefinitions.apikey TokenAuth
// @in header
// @name Authorization
func main() {
	// Parse arguments
	configPath := flag.String("config", config.DefaultConfigPath, "Config file path")

	// Load config
	cfg, err := loadConfig(*configPath)
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Init auth service
	authService, err := initAuthService(cfg)
	if err != nil {
		panic(fmt.Sprintf("Failed to init auth service: %v", err))
	}

	// Init DB
	mongoDB, err := initMongoDB(cfg)
	if err != nil {
		panic(fmt.Sprintf("Failed to init MongoDB: %v", err))
	}

	// Init cache
	redisCache, err := initRedisCache(cfg)
	if err != nil {
		panic(fmt.Sprintf("Failed to init Redis cache: %v", err))
	}

	// Init repositories
	repositories := initRepositories(mongoDB, redisCache, cfg)

	// Setup server
	engine := gin.Default()
	registerAPIRouters(engine, authService, repositories)

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
func loadConfig(path string) (*config.Config, error) {
	cfg, err := config.LoadConfig(path)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// Init auth service
func initAuthService(cfg *config.Config) (coreAuth.BaseAuthService, error) {
	authService, err := coreAuth.NewAuthService(context.Background(), &cfg.AuthService)
	if err != nil {
		return nil, err
	}

	return authService, nil
}

// Init MongoDB
func initMongoDB(cfg *config.Config) (*coreDB.MongoDB, error) {
	mongoDB, err := coreDB.NewMongoDB(context.Background(), &cfg.MongoDB)
	if err != nil {
		return nil, err
	}

	return mongoDB, nil
}

// Init Redis cache
func initRedisCache(cfg *config.Config) (*coreCache.RedisCache, error) {
	redisCache, err := coreCache.NewRedisCache(context.Background(), &cfg.RedisCache)
	if err != nil {
		return nil, err
	}

	return redisCache, nil
}

// Init repositories
func initRepositories(mongoDB *coreDB.MongoDB, redisCache *coreCache.RedisCache, cfg *config.Config) map[coreRepository.RepositoryName]any {
	repositories := make(map[coreRepository.RepositoryName]any)

	// Init user db repository
	userDBRepo := coreRepository.NewUserMongoDBRepository(mongoDB, cfg.Repositories.UserDB)
	repositories[coreRepository.RepositoryNameUserDB] = userDBRepo

	// Init user cache repository
	userCacheRepo := coreRepository.NewUserRedisCacheRepository(redisCache, cfg.Repositories.UserCache)
	repositories[coreRepository.RepositoryNameUserCache] = userCacheRepo

	return repositories
}

// Register API routers
func registerAPIRouters(engine *gin.Engine, authService coreAuth.BaseAuthService, repositories map[coreRepository.RepositoryName]any) {
	userDBRepo, _ := repositories[coreRepository.RepositoryNameUserDB].(coreRepository.UserDBRepository)
	userCacheRepo, _ := repositories[coreRepository.RepositoryNameUserCache].(coreRepository.UserCacheRepository)

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
	v1RouterGroup.Use(middleware.TokenAuthenticationHandler(authService, userDBRepo, userCacheRepo))

	// Register v1 user router
	userRouterGroup := v1RouterGroup.Group("/user")
	router.RegisterV1UserRouter(userRouterGroup)
}
