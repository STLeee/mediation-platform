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

// @Summary Liveness check
// @Description Liveness check
// @Tags health
// @Router /health/liveness [get]
// @Produce json
// @Success 200 {object} model.MessageResponse
func (hc *HealthController) Liveness(c *gin.Context) {
	hc.ResponseOK(c)
}

// @Summary Readiness check
// @Description Readiness check
// @Tags health
// @Router /health/readiness [get]
// @Produce json
// @Success 200 {object} model.MessageResponse
func (hc *HealthController) Readiness(c *gin.Context) {
	hc.ResponseOK(c)
}
