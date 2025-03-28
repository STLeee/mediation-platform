package v1

import (
	"net/http"

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
// @Router /v1/user/{user_id} [get]
// @Security TokenAuth
// @Param user_id path string true "User ID"
// @Produce json
// @Success 200 {object} model.GetUserResponse
func (hc *UserController) GetUser(c *gin.Context) {
	user := c.MustGet("user").(*coreModel.User)
	userID := c.Param("user_id")

	if userID != user.UserID {
		c.Error(model.HttpStatusCodeError{
			StatusCode: http.StatusForbidden,
			Message:    "User ID does not match",
		})
		c.Abort()
		return
	}

	c.JSON(200, model.GetUserResponse{
		UserID:      user.UserID,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		PhotoURL:    user.PhotoURL,
	})
}
