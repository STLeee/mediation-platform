package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/STLeee/mediation-platform/backend/app/api-service/model"
	coreAuth "github.com/STLeee/mediation-platform/backend/core/auth"
	coreRepository "github.com/STLeee/mediation-platform/backend/core/repository"
)

// TokenAuthenticationHandler is a middleware for token authentication
func TokenAuthenticationHandler(authService coreAuth.BaseAuthService, userRepo coreRepository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from header
		token := c.GetHeader("Authorization")
		if len(token) < 8 || token[:7] != "Bearer " {
			c.Next()
			return
		}
		token = token[7:]

		// Authenticate user by token
		authUID, err := authService.AuthenticateByToken(c, token)
		if err != nil {
			responseError := model.HttpStatusCodeError{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed to authenticate user by token",
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

		// Get user from MongoDB
		user, err := userRepo.GetUserByAuthUID(c, authService.GetName(), authUID)
		if err != nil {
			// If user not found, create a new user
			if repositoryError, ok := err.(coreRepository.RepositoryError); ok && repositoryError.ErrType == coreRepository.RepositoryErrorTypeRecordNotFound {
				userFromAuth, err := authService.GetUserInfo(c, authUID)
				if err != nil {
					responseError := model.HttpStatusCodeError{
						StatusCode: http.StatusInternalServerError,
						Message:    "failed to get user info from auth service",
						Err:        err,
					}
					c.Error(responseError)
					c.Abort()
					return
				}
				userID, err := userRepo.CreateUser(c, userFromAuth)
				if err != nil {
					responseError := model.HttpStatusCodeError{
						StatusCode: http.StatusInternalServerError,
						Message:    "failed to create user",
						Err:        err,
					}
					c.Error(responseError)
					c.Abort()
					return
				}
				user, err = userRepo.GetUserByID(c, userID)
				if err != nil {
					responseError := model.HttpStatusCodeError{
						StatusCode: http.StatusInternalServerError,
						Message:    "failed to get user",
						Err:        err,
					}
					c.Error(responseError)
					c.Abort()
					return
				}
			} else {
				responseError := model.HttpStatusCodeError{
					StatusCode: http.StatusInternalServerError,
					Message:    "failed to get user by auth UID",
					Err:        err,
				}
				c.Error(responseError)
				c.Abort()
				return
			}
		}

		// Set user info to context
		c.Set("user", user)
		c.Next()
	}
}
