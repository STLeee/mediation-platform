package router

import (
	"github.com/gin-gonic/gin"

	"github.com/STLeee/mediation-platform/backend/app/api-service/controller"
	controllerV1 "github.com/STLeee/mediation-platform/backend/app/api-service/controller/v1"
	middlewareV1 "github.com/STLeee/mediation-platform/backend/app/api-service/middleware/v1"
)

func RegisterHealthRouter(r *gin.RouterGroup) {
	healthController := controller.NewHealthController()

	r.GET("/liveness", healthController.Liveness)
	r.GET("/readiness", healthController.Readiness)
}

func RegisterV1UserRouter(r *gin.RouterGroup) {
	r.Use(middlewareV1.UserAPIAuthorizationHandler())

	userController := controllerV1.NewUserController()

	r.GET("/:user_id", userController.GetUser)
}
