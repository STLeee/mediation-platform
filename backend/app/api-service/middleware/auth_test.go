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
	AuthenticateByTokenFunc   func(ctx context.Context, token string) (uid string, err error)
	GetUserInfoAndMappingFunc func(ctx context.Context, uid string) (user *coreModel.User, mapping map[string]any, err error)
}

func (authService *TestAuthService) AuthenticateByToken(ctx context.Context, token string) (uid string, err error) {
	return authService.AuthenticateByTokenFunc(ctx, token)
}

func (authService *TestAuthService) GetUserInfoAndMapping(ctx context.Context, uid string) (user *coreModel.User, mapping map[string]any, err error) {
	return authService.GetUserInfoAndMappingFunc(ctx, uid)
}

func TestTokenAuthenticationHandler(t *testing.T) {
	testCases := []struct {
		name                   string
		user                   *coreModel.User
		authenticateByTokenErr error
		getUserInfoFuncErr     error
		expected_code          int
	}{
		{
			name:          "no-error",
			user:          &coreModel.User{UserID: "test_user_id"},
			expected_code: http.StatusOK,
		},
		{
			name:                   "authenticate-by-token/server-error",
			user:                   nil,
			authenticateByTokenErr: coreAuth.AuthServiceError{ErrType: coreAuth.AuthServiceErrorTypeServerError},
			expected_code:          http.StatusInternalServerError,
		},
		{
			name:                   "authenticate-by-token-error/token-invalid",
			user:                   nil,
			authenticateByTokenErr: coreAuth.AuthServiceError{ErrType: coreAuth.AuthServiceErrorTypeTokenInvalid},
			expected_code:          http.StatusUnauthorized,
		},
		{
			name:                   "authenticate-by-token-error/user-not-found",
			user:                   nil,
			authenticateByTokenErr: coreAuth.AuthServiceError{ErrType: coreAuth.AuthServiceErrorTypeUserNotFound},
			expected_code:          http.StatusUnauthorized,
		},
		{
			name:               "get-user-info-error/server-error",
			user:               nil,
			getUserInfoFuncErr: coreAuth.AuthServiceError{ErrType: coreAuth.AuthServiceErrorTypeServerError},
			expected_code:      http.StatusInternalServerError,
		},
		{
			name:               "get-user-info-error/user-not-found",
			user:               nil,
			getUserInfoFuncErr: coreAuth.AuthServiceError{ErrType: coreAuth.AuthServiceErrorTypeUserNotFound},
			expected_code:      http.StatusUnauthorized,
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
				GetUserInfoAndMappingFunc: func(ctx context.Context, uid string) (user *coreModel.User, mapping map[string]any, err error) {
					if testCase.getUserInfoFuncErr != nil {
						return nil, nil, testCase.getUserInfoFuncErr
					}
					return testCase.user, map[string]any{"test_auth_uid": "test_user_id"}, nil
				},
			}
			httpRecorder := utils.RegisterAndRecordHttpRequest(func(routeGroup *gin.RouterGroup) {
				routeGroup.Use(ErrorHandler(), TokenAuthenticationHandler(testAuthService, nil))
				routeGroup.Handle("GET", "/test", func(c *gin.Context) {
					c.JSON(200, c.MustGet("user"))
				})
			}, "GET", "/test", nil)

			assert.Equal(t, testCase.expected_code, httpRecorder.Code)
			if httpRecorder.Code == 200 {
				assert.Equal(t, utils.ConvertToJSONString(testCase.user), httpRecorder.Body.String())
			}
		})
	}
}
