package middleware

import (
	"context"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	coreAuth "github.com/STLeee/mediation-platform/backend/core/auth"
	coreModel "github.com/STLeee/mediation-platform/backend/core/model"
	"github.com/STLeee/mediation-platform/backend/core/utils"
)

type TestAuthService struct {
	AuthenticateByTokenFunc func(ctx context.Context, token string) (uid string, err error)
	GetUserInfoFunc         func(ctx context.Context, uid string) (*coreModel.UserInfo, error)
}

func (authService *TestAuthService) AuthenticateByToken(ctx context.Context, token string) (uid string, err error) {
	return authService.AuthenticateByTokenFunc(ctx, token)
}

func (authService *TestAuthService) GetUserInfo(ctx context.Context, uid string) (*coreModel.UserInfo, error) {
	return authService.GetUserInfoFunc(ctx, uid)
}

func TestTokenAuthHandler(t *testing.T) {
	testCases := []struct {
		name                   string
		user                   *coreModel.UserInfo
		authenticateByTokenErr error
		getUserInfoFuncErr     error
		excepted_code          int
	}{
		{
			name:          "no-error",
			user:          &coreModel.UserInfo{UserID: "test_user_id"},
			excepted_code: http.StatusOK,
		},
		{
			name:                   "authenticate-by-token/server-error",
			user:                   nil,
			authenticateByTokenErr: coreAuth.AuthServiceError{ErrType: coreAuth.AuthServiceErrorTypeServerError},
			excepted_code:          http.StatusInternalServerError,
		},
		{
			name:                   "authenticate-by-token-error/token-invalid",
			user:                   nil,
			authenticateByTokenErr: coreAuth.AuthServiceError{ErrType: coreAuth.AuthServiceErrorTypeTokenInvalid},
			excepted_code:          http.StatusUnauthorized,
		},
		{
			name:                   "authenticate-by-token-error/user-not-found",
			user:                   nil,
			authenticateByTokenErr: coreAuth.AuthServiceError{ErrType: coreAuth.AuthServiceErrorTypeUserNotFound},
			excepted_code:          http.StatusUnauthorized,
		},
		{
			name:               "get-user-info-error/server-error",
			user:               nil,
			getUserInfoFuncErr: coreAuth.AuthServiceError{ErrType: coreAuth.AuthServiceErrorTypeServerError},
			excepted_code:      http.StatusInternalServerError,
		},
		{
			name:               "get-user-info-error/user-not-found",
			user:               nil,
			getUserInfoFuncErr: coreAuth.AuthServiceError{ErrType: coreAuth.AuthServiceErrorTypeUserNotFound},
			excepted_code:      http.StatusUnauthorized,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testAuthService := &TestAuthService{
				AuthenticateByTokenFunc: func(ctx context.Context, token string) (uid string, err error) {
					if testCase.authenticateByTokenErr != nil {
						return "", testCase.authenticateByTokenErr
					}
					return "test_user_id", nil
				},
				GetUserInfoFunc: func(ctx context.Context, uid string) (*coreModel.UserInfo, error) {
					if testCase.getUserInfoFuncErr != nil {
						return nil, testCase.getUserInfoFuncErr
					}
					return testCase.user, nil
				},
			}
			httpRecorder := utils.RegisterAndRecordHttpRequest(func(routeGroup *gin.RouterGroup) {
				routeGroup.Use(ErrorHandler(), TokenAuthHandler(testAuthService))
				routeGroup.Handle("GET", "/test", func(c *gin.Context) {
					c.JSON(200, c.MustGet("user"))
				})
			}, "GET", "/test", nil)

			assert.Equal(t, testCase.excepted_code, httpRecorder.Code)
			if httpRecorder.Code == 200 {
				assert.Equal(t, utils.ToJSONString(testCase.user), httpRecorder.Body.String())
			}
		})
	}
}
