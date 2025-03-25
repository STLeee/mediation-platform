package middleware

import (
	"context"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	coreAuth "github.com/STLeee/mediation-platform/backend/core/auth"
	coreModel "github.com/STLeee/mediation-platform/backend/core/model"
	coreRepository "github.com/STLeee/mediation-platform/backend/core/repository"
	"github.com/STLeee/mediation-platform/backend/core/utils"
)

type TestUserRepository struct {
	GetUserByTokenFunc func(ctx context.Context, token string) (*coreModel.User, error)
}

func (repo *TestUserRepository) GetUserByToken(ctx context.Context, token string) (*coreModel.User, error) {
	return repo.GetUserByTokenFunc(ctx, token)
}

func (repo *TestUserRepository) GetUserByID(ctx context.Context, userID string) (*coreModel.User, error) {
	return nil, nil
}

func TestTokenAuthenticationHandler(t *testing.T) {
	testCases := []struct {
		name           string
		user           *coreModel.User
		getUserByIDErr error
		expected_code  int
	}{
		{
			name:          "no-error",
			user:          &coreModel.User{UserID: "test_user_id"},
			expected_code: http.StatusOK,
		},
		{
			name:           "auth/server-error",
			user:           nil,
			getUserByIDErr: coreAuth.AuthServiceError{ErrType: coreAuth.AuthServiceErrorTypeServerError},
			expected_code:  http.StatusInternalServerError,
		},
		{
			name:           "auth/token-invalid",
			user:           nil,
			getUserByIDErr: coreAuth.AuthServiceError{ErrType: coreAuth.AuthServiceErrorTypeTokenInvalid},
			expected_code:  http.StatusUnauthorized,
		},
		{
			name:           "auth/user-not-found",
			user:           nil,
			getUserByIDErr: coreAuth.AuthServiceError{ErrType: coreAuth.AuthServiceErrorTypeUserNotFound},
			expected_code:  http.StatusUnauthorized,
		},
		{
			name:           "db/server-error",
			user:           nil,
			getUserByIDErr: coreRepository.RepositoryError{ErrType: coreRepository.RepositoryErrorTypeServerError},
			expected_code:  http.StatusInternalServerError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testUserRepo := &TestUserRepository{
				GetUserByTokenFunc: func(ctx context.Context, token string) (*coreModel.User, error) {
					return testCase.user, testCase.getUserByIDErr
				},
			}
			httpRecorder := utils.RegisterAndRecordHttpRequest(func(routeGroup *gin.RouterGroup) {
				routeGroup.Use(ErrorHandler(), TokenAuthenticationHandler(testUserRepo))
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
