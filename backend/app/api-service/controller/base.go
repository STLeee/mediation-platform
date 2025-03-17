package controller

import (
	"github.com/gin-gonic/gin"

	controllerModel "github.com/STLeee/mediation-platform/backend/app/api-service/model/controller"
)

// BaseController is a base controller
type BaseController struct{}

// ResponseOK is a response for ok
func (bc *BaseController) ResponseOK(c *gin.Context) {
	c.JSON(200, controllerModel.NewMessageResponse("ok"))
}
