package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/STLeee/mediation-platform/backend/app/api-service/model"
	coreModel "github.com/STLeee/mediation-platform/backend/core/model"
)

// UserAPIAuthorizationHandler is a middleware for User API authorization
func UserAPIAuthorizationHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context
		var user *coreModel.User
		if userInterface, ok := c.Get("user"); ok {
			user, _ = userInterface.(*coreModel.User)
		}
		if user == nil {
			c.Error(model.HttpStatusCodeError{
				StatusCode: http.StatusInternalServerError,
				Message:    "User not found in context",
			})
			c.Abort()
			return
		}

		// Check is user is owner of the resource
		userID := c.Param("user_id")
		if userID == user.UserID {
			c.Next()
			return
		}

		c.Error(model.HttpStatusCodeError{
			StatusCode: http.StatusForbidden,
		})
		c.Abort()
	}
}
