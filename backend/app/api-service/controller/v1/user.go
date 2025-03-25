package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/STLeee/mediation-platform/backend/app/api-service/model"
	coreModel "github.com/STLeee/mediation-platform/backend/core/model"
)

// UserController is a controller for user management
type UserController struct{}

// NewUserController creates a new UserController
func NewUserController() *UserController {
	return &UserController{}
}

// @Summary Get user
// @Description Get user info
// @Tags user
// @Router /v1/user/:user_id [get]
// @Security TokenAuth
// @Param user_id path string true "User ID"
// @Produce json
// @Success 200 {object} model.GetUserResponse
func (hc *UserController) GetUser(c *gin.Context) {
	userInfo := c.MustGet("user").(*coreModel.User)
	c.JSON(200, model.GetUserResponse{
		UserID:      userInfo.UserID,
		DisplayName: userInfo.DisplayName,
		Email:       userInfo.Email,
		PhoneNumber: userInfo.PhoneNumber,
		PhotoURL:    userInfo.PhotoURL,
	})
}
