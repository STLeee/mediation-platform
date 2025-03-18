package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/STLeee/mediation-platform/backend/app/api-service/model"
)

// BaseController is a base controller
type BaseController struct{}

// ResponseOK is a response for ok
func (bc *BaseController) ResponseOK(c *gin.Context) {
	c.JSON(200, model.NewMessageResponse("ok"))
}
