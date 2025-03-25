package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/STLeee/mediation-platform/backend/app/api-service/model"
	coreAuth "github.com/STLeee/mediation-platform/backend/core/auth"
	coreRepository "github.com/STLeee/mediation-platform/backend/core/repository"
)

// TokenAuthenticationHandler is a middleware for token authentication
func TokenAuthenticationHandler(userRepo coreRepository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from header
		token := c.GetHeader("Authorization")

		// Authenticate user by token
		user, err := userRepo.GetUserByToken(context.Background(), token)
		if err != nil {
			responseError := model.HttpStatusCodeError{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
			if authServiceError, ok := err.(coreAuth.AuthServiceError); ok {
				errType := authServiceError.ErrType
				if errType == coreAuth.AuthServiceErrorTypeTokenInvalid || errType == coreAuth.AuthServiceErrorTypeUserNotFound {
					responseError = model.HttpStatusCodeError{
						StatusCode: http.StatusUnauthorized,
						Err:        err,
					}
				}
			}
			c.Error(responseError)
			c.Abort()
			return
		}

		// Set user info to context
		c.Set("user", user)
		c.Next()
	}
}
