package controller

import (
	"github.com/gin-gonic/gin"
)

// HealthController is a controller for health check
type HealthController struct {
	BaseController
}

// NewHealthController creates a new HealthController
func NewHealthController() *HealthController {
	return &HealthController{}
}

// Liveness is a handler for liveness check
func (hc *HealthController) Liveness(c *gin.Context) {
	hc.ResponseOK(c)
}

// Readiness is a handler for readiness check
func (hc *HealthController) Readiness(c *gin.Context) {
	hc.ResponseOK(c)
}
