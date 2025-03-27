package middleware

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	coreAuth "github.com/STLeee/mediation-platform/backend/core/auth"
	coreModel "github.com/STLeee/mediation-platform/backend/core/model"
	coreRepository "github.com/STLeee/mediation-platform/backend/core/repository"
	"github.com/STLeee/mediation-platform/backend/core/utils"
)

type MockFirebaseAuthService struct {
	AuthenticateByTokenFunc func(ctx context.Context, token string) (uid string, err error)
	GetUserInfoFunc         func(ctx context.Context, uid string) (user *coreModel.User, err error)
}

func (auth *MockFirebaseAuthService) GetName() coreAuth.AuthServiceName {
	return coreAuth.AuthServiceNameFirebase
}

func (auth *MockFirebaseAuthService) AuthenticateByToken(ctx context.Context, token string) (uid string, err error) {
	return auth.AuthenticateByTokenFunc(ctx, token)
}

func (auth *MockFirebaseAuthService) GetUserInfo(ctx context.Context, uid string) (user *coreModel.User, err error) {
	return auth.GetUserInfoFunc(ctx, uid)
}

type MockUserDBRepository struct {
	CreateUserFunc       func(ctx context.Context, user *coreModel.User) (string, error)
	GetUserByAuthUIDFunc func(ctx context.Context, authName coreAuth.AuthServiceName, authUID string) (*coreModel.User, error)
	GetUserByIDFunc      func(ctx context.Context, userID string) (*coreModel.User, error)
}

func (repo *MockUserDBRepository) CreateUser(ctx context.Context, user *coreModel.User) (string, error) {
	return repo.CreateUserFunc(ctx, user)
}

func (repo *MockUserDBRepository) GetUserByAuthUID(ctx context.Context, authName coreAuth.AuthServiceName, authUID string) (*coreModel.User, error) {
	return repo.GetUserByAuthUIDFunc(ctx, authName, authUID)
}

func (repo *MockUserDBRepository) GetUserByID(ctx context.Context, userID string) (*coreModel.User, error) {
	return repo.GetUserByIDFunc(ctx, userID)
}

var mockFirebaseUser = &coreModel.User{
	UserID:      "",
	FirebaseUID: "test-firebase-uid",
	DisplayName: "Test User",
	Email:       "test-user@mediation-platform.com",
	PhoneNumber: "",
	PhotoURL:    "",
	Disabled:    false,
}

var mockUserInDB = &coreModel.User{
	UserID:      "test-user-id",
	FirebaseUID: "test-firebase-uid",
	DisplayName: "Test User",
	Email:       "test-user@mediation-platform.com",
	PhoneNumber: "",
	PhotoURL:    "",
	Disabled:    false,
	CreatedAt:   time.Now(),
	UpdatedAt:   time.Now(),
	LastLoginAt: time.Now(),
}

func TestTokenAuthenticationHandler(t *testing.T) {
	testCases := []struct {
		name                       string
		token                      string
		authUID                    string
		authUser                   *coreModel.User
		dbUser                     *coreModel.User
		authenticateByTokenFuncErr error
		getUserByAuthUIDFuncErr    error
		getUserInfoFuncErr         error
		createUserFuncErr          error
		getUserByIDErr             error
		expectedStatusCode         int
	}{
		{
			name:               "success",
			token:              "test-token",
			authUID:            mockFirebaseUser.FirebaseUID,
			authUser:           mockFirebaseUser,
			dbUser:             mockUserInDB,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "empty-token",
			token:              "",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:                       "auth/invalid-token",
			token:                      "invalid-token",
			authenticateByTokenFuncErr: coreAuth.AuthServiceError{ErrType: coreAuth.AuthServiceErrorTypeTokenInvalid},
			expectedStatusCode:         http.StatusUnauthorized,
		},
		{
			name:                       "auth/user-not-found",
			token:                      "test-token",
			authenticateByTokenFuncErr: coreAuth.AuthServiceError{ErrType: coreAuth.AuthServiceErrorTypeUserNotFound},
			expectedStatusCode:         http.StatusUnauthorized,
		},
		{
			name:                       "auth/unknown-error",
			token:                      "test-token",
			authenticateByTokenFuncErr: coreAuth.AuthServiceError{ErrType: coreAuth.AuthServiceErrorTypeServerError},
			expectedStatusCode:         http.StatusInternalServerError,
		},
		{
			name:                    "db/get-user-by-auth-uid-error",
			token:                   "test-token",
			authUID:                 mockFirebaseUser.FirebaseUID,
			authUser:                mockFirebaseUser,
			getUserByAuthUIDFuncErr: coreRepository.RepositoryError{ErrType: coreRepository.RepositoryErrorTypeServerError},
			expectedStatusCode:      http.StatusInternalServerError,
		},
		{
			name:                    "auth/get-user-info-error",
			token:                   "test-token",
			authUID:                 mockFirebaseUser.FirebaseUID,
			getUserByAuthUIDFuncErr: coreRepository.RepositoryError{ErrType: coreRepository.RepositoryErrorTypeRecordNotFound},
			getUserInfoFuncErr:      coreAuth.AuthServiceError{ErrType: coreAuth.AuthServiceErrorTypeServerError},
			expectedStatusCode:      http.StatusInternalServerError,
		},
		{
			name:                    "db/create-user-error",
			token:                   "test-token",
			authUID:                 mockFirebaseUser.FirebaseUID,
			authUser:                mockFirebaseUser,
			getUserByAuthUIDFuncErr: coreRepository.RepositoryError{ErrType: coreRepository.RepositoryErrorTypeRecordNotFound},
			createUserFuncErr:       coreAuth.AuthServiceError{ErrType: coreAuth.AuthServiceErrorTypeServerError},
			expectedStatusCode:      http.StatusInternalServerError,
		},
		{
			name:                    "db/get-user-by-id-error",
			token:                   "test-token",
			authUID:                 mockFirebaseUser.FirebaseUID,
			authUser:                mockFirebaseUser,
			dbUser:                  mockUserInDB,
			getUserByAuthUIDFuncErr: coreRepository.RepositoryError{ErrType: coreRepository.RepositoryErrorTypeRecordNotFound},
			getUserByIDErr:          coreRepository.RepositoryError{ErrType: coreRepository.RepositoryErrorTypeServerError},
			expectedStatusCode:      http.StatusInternalServerError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockFirebaseAuthService := &MockFirebaseAuthService{
				AuthenticateByTokenFunc: func(ctx context.Context, token string) (uid string, err error) {
					assert.Equal(t, testCase.token, token)
					return testCase.authUID, testCase.authenticateByTokenFuncErr
				},
				GetUserInfoFunc: func(ctx context.Context, uid string) (user *coreModel.User, err error) {
					assert.Equal(t, testCase.authUID, uid)
					return testCase.authUser, testCase.getUserInfoFuncErr
				},
			}
			mockUserDBRepo := &MockUserDBRepository{
				CreateUserFunc: func(ctx context.Context, user *coreModel.User) (string, error) {
					assert.Equal(t, testCase.authUser, user)
					userID := ""
					if testCase.dbUser != nil {
						userID = testCase.dbUser.UserID
					}
					return userID, testCase.createUserFuncErr
				},
				GetUserByAuthUIDFunc: func(ctx context.Context, authName coreAuth.AuthServiceName, authUID string) (*coreModel.User, error) {
					assert.Equal(t, coreAuth.AuthServiceNameFirebase, authName)
					assert.Equal(t, testCase.authUID, authUID)
					return testCase.dbUser, testCase.getUserByAuthUIDFuncErr
				},
				GetUserByIDFunc: func(ctx context.Context, userID string) (*coreModel.User, error) {
					assert.Equal(t, testCase.dbUser.UserID, userID)
					return testCase.dbUser, testCase.getUserByIDErr
				},
			}

			httpRecorder := utils.RegisterAndRecordHttpRequest(func(routeGroup *gin.RouterGroup) {
				routeGroup.Use(func(ctx *gin.Context) {
					ctx.Request.Header.Set("Authorization", "Bearer "+testCase.token)
					ctx.Next()
				})
				routeGroup.Use(ErrorHandler(), TokenAuthenticationHandler(mockFirebaseAuthService, mockUserDBRepo, nil))
				routeGroup.Handle("GET", "/test", func(c *gin.Context) {
					if testCase.token == "" {
						c.JSON(http.StatusUnauthorized, nil)
						return
					}
					c.JSON(http.StatusOK, c.MustGet("user"))
				})
			}, "GET", "/test", nil)

			assert.Equal(t, testCase.expectedStatusCode, httpRecorder.Code)
			if httpRecorder.Code == http.StatusOK {
				assert.Equal(t, utils.ConvertToJSONString(testCase.dbUser), httpRecorder.Body.String())
			}
		})
	}
}
