package router

import (
	"github.com/gin-gonic/gin"

	"github.com/STLeee/mediation-platform/backend/app/api-service/controller"
)

func RegisterHealthRouter(r *gin.RouterGroup) {
	healthController := controller.NewHealthController()

	r.GET("/liveness", healthController.Liveness)
	r.GET("/readiness", healthController.Readiness)
}
