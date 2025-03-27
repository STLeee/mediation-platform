package v1

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/STLeee/mediation-platform/backend/app/api-service/model"
	coreModel "github.com/STLeee/mediation-platform/backend/core/model"
	"github.com/STLeee/mediation-platform/backend/core/utils"
)

func TestUserAPIAuthorizationHandler(t *testing.T) {
	testCases := []struct {
		name       string
		tokenUser  *coreModel.User
		queryPath  string
		statusCode int
		expected   *model.GetUserResponse
		isErr      bool
	}{
		{
			name: "user-owner",
			tokenUser: &coreModel.User{
				UserID:      "test-user-id",
				DisplayName: "test-display-name",
				Email:       "test-email",
				PhoneNumber: "test-phone-number",
				PhotoURL:    "test-photo-url",
			},
			queryPath:  "/test-user-id",
			statusCode: http.StatusOK,
			expected: &model.GetUserResponse{
				UserID:      "test-user-id",
				DisplayName: "test-display-name",
				Email:       "test-email",
				PhoneNumber: "test-phone-number",
				PhotoURL:    "test-photo-url",
			},
		},
		{
			name: "user-not-owner",
			tokenUser: &coreModel.User{
				UserID:      "test-user-id",
				DisplayName: "test-display-name",
				Email:       "test-email",
				PhoneNumber: "test-phone-number",
				PhotoURL:    "test-photo-url",
			},
			queryPath:  "/test-user-id-2",
			statusCode: http.StatusForbidden,
			isErr:      true,
		},
		{
			name:       "no-user",
			tokenUser:  nil,
			queryPath:  "/test-user-id",
			statusCode: http.StatusUnauthorized,
			isErr:      true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			httpRecorder := utils.RegisterAndRecordHttpRequest(
				func(router *gin.RouterGroup) {
					router.Use(func(ctx *gin.Context) {
						// Set user to context
						ctx.Set("user", testCase.tokenUser)
						ctx.Next()

						// Check error
						if testCase.isErr {
							err := ctx.Errors.Last()
							assert.NotNil(t, err)
							assert.Equal(t, testCase.statusCode, err.Err.(model.HttpStatusCodeError).StatusCode)
						}
					})
					router.Use(UserAPIAuthorizationHandler())
					router.GET("/:user_id", func(ctx *gin.Context) {
						ctx.JSON(http.StatusOK, testCase.expected)
					})
				},
				"GET",
				testCase.queryPath,
				nil,
			)

			// Check response
			if !testCase.isErr {
				assert.Equal(t, testCase.statusCode, httpRecorder.Code)
				assert.Equal(t, utils.ConvertToJSONString(testCase.expected), httpRecorder.Body.String())
			}
		})
	}
}
