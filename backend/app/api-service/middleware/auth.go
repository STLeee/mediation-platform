package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/STLeee/mediation-platform/backend/app/api-service/model"
	coreAuth "github.com/STLeee/mediation-platform/backend/core/auth"
)

// TokenAuthHandler is a middleware for token authentication
func TokenAuthHandler(authService coreAuth.BaseAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from header
		token := c.GetHeader("Authorization")

		// Authenticate user by token
		uid, err := authService.AuthenticateByToken(c, token)
		if err != nil {
			responseError := model.HttpStatusCodeError{StatusCode: http.StatusInternalServerError}
			if errors.As(err, &coreAuth.AuthServiceError{}) {
				errType := err.(coreAuth.AuthServiceError).ErrType
				if errType == coreAuth.AuthServiceErrorTypeTokenInvalid || errType == coreAuth.AuthServiceErrorTypeUserNotFound {
					responseError = model.HttpStatusCodeError{StatusCode: http.StatusUnauthorized}
				}
			}
			c.Error(responseError)
			c.Abort()
			return
		}

		// Get user info from auth service
		userInfo, err := authService.GetUserInfo(c, uid)
		if err != nil {
			responseError := model.HttpStatusCodeError{StatusCode: http.StatusInternalServerError}
			if errors.As(err, &coreAuth.AuthServiceError{}) {
				errType := err.(coreAuth.AuthServiceError).ErrType
				if errType == coreAuth.AuthServiceErrorTypeUserNotFound {
					responseError = model.HttpStatusCodeError{StatusCode: http.StatusUnauthorized}
				}
			}
			c.Error(responseError)
			c.Abort()
			return
		}

		// Set user info to context
		c.Set("user", userInfo)
		c.Next()
	}
}
