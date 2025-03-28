package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/STLeee/mediation-platform/backend/app/api-service/model"
	coreAuth "github.com/STLeee/mediation-platform/backend/core/auth"
	coreModel "github.com/STLeee/mediation-platform/backend/core/model"
	coreRepository "github.com/STLeee/mediation-platform/backend/core/repository"
)

func authenticateUserByToken(ctx context.Context, authService coreAuth.BaseAuthService, userDBRepo coreRepository.UserDBRepository, token string) (*coreModel.User, error) {
	// Empty token
	if token == "" {
		return nil, model.HttpStatusCodeError{
			StatusCode: http.StatusUnauthorized,
			Message:    "empty token",
		}
	}

	// Authenticate user by token
	authUID, err := authService.AuthenticateByToken(ctx, token)
	if err != nil {
		if authServiceError, ok := err.(coreAuth.AuthServiceError); ok {
			errType := authServiceError.ErrType
			if errType == coreAuth.AuthServiceErrorTypeTokenInvalid || errType == coreAuth.AuthServiceErrorTypeUserNotFound {
				return nil, model.HttpStatusCodeError{
					StatusCode: http.StatusUnauthorized,
					Err:        err,
				}
			}
		}
		return nil, model.HttpStatusCodeError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to authenticate user by token",
			Err:        err,
		}
	}

	// Get user from MongoDB
	user, err := userDBRepo.GetUserByAuthUID(ctx, authService.GetName(), authUID)
	if err != nil {
		// If user not found, create a new user
		if repositoryError, ok := err.(coreRepository.RepositoryError); ok && repositoryError.ErrType == coreRepository.RepositoryErrorTypeRecordNotFound {
			userFromAuth, err := authService.GetUserInfo(ctx, authUID)
			if err != nil {
				return nil, model.HttpStatusCodeError{
					StatusCode: http.StatusInternalServerError,
					Message:    "failed to get user info from auth service",
					Err:        err,
				}
			}
			userID, err := userDBRepo.CreateUser(ctx, userFromAuth)
			if err != nil {
				return nil, model.HttpStatusCodeError{
					StatusCode: http.StatusInternalServerError,
					Message:    "failed to create user",
					Err:        err,
				}
			}
			user, err = userDBRepo.GetUserByID(ctx, userID)
			if err != nil {
				return nil, model.HttpStatusCodeError{
					StatusCode: http.StatusInternalServerError,
					Message:    "failed to get user",
					Err:        err,
				}
			}
		} else {
			return nil, model.HttpStatusCodeError{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed to get user by auth UID",
				Err:        err,
			}
		}
	}

	return user, nil
}

func newErrorUser(err model.HttpStatusCodeError) *coreModel.User {
	return &coreModel.User{
		UserID:      "error",
		DisplayName: err.Message,
		PhoneNumber: strconv.Itoa(err.StatusCode),
		Disabled:    true,
	}
}

func parseErrorUser(user *coreModel.User) (*coreModel.User, error) {
	if user == nil {
		return nil, nil
	}
	if user.UserID == "error" {
		statusCode, _ := strconv.Atoi(user.PhoneNumber)
		return nil, model.HttpStatusCodeError{
			StatusCode: statusCode,
			Message:    user.DisplayName,
		}
	}
	return user, nil
}

// TokenAuthenticationHandler is a middleware for token authentication
func TokenAuthenticationHandler(authService coreAuth.BaseAuthService, userDBRepo coreRepository.UserDBRepository, userCacheRepo coreRepository.UserCacheRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from header
		token := c.GetHeader("Authorization")
		if len(token) < 8 || token[:7] != "Bearer " {
			c.Next()
			return
		}
		token = token[7:]

		var user *coreModel.User

		// Get user from cache
		if userCacheRepo != nil {
			var cacheErr error
			user, cacheErr = userCacheRepo.GetAuthTokenUser(c, authService.GetName(), token)
			if cacheErr != nil {
				if repositoryError, ok := cacheErr.(coreRepository.RepositoryError); ok && repositoryError.ErrType == coreRepository.RepositoryErrorTypeRecordNotFound {
					user = nil
				} else {
					// TODO: record error
				}
			} else {
				var err error
				user, err = parseErrorUser(user)
				if err != nil {
					c.Error(err)
					c.Abort()
					return
				}
			}
		}

		if user == nil {
			// Authenticate user by token
			var err error
			user, err = authenticateUserByToken(c, authService, userDBRepo, token)
			if err != nil {
				if httpStatusCodeError, ok := err.(model.HttpStatusCodeError); ok && httpStatusCodeError.StatusCode != http.StatusInternalServerError {
					// Set error to cache
					if userCacheRepo != nil {
						errUser := newErrorUser(httpStatusCodeError)
						if cacheErr := userCacheRepo.SetAuthTokenUser(c, authService.GetName(), token, errUser); cacheErr != nil {
							// TODO: record error
						}
					}
				}
				c.Error(err)
				c.Abort()
				return
			}

			// Set user to cache
			if userCacheRepo != nil {
				err = userCacheRepo.SetAuthTokenUser(c, authService.GetName(), token, user)
				if err != nil {
					// TODO: record error
				}
			}
		}

		// Set user info to context
		c.Set("user", user)
		c.Next()
	}
}
