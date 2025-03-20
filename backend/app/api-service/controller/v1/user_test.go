package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/STLeee/mediation-platform/backend/app/api-service/model"
	coreModel "github.com/STLeee/mediation-platform/backend/core/model"
	"github.com/STLeee/mediation-platform/backend/core/utils"
)

func TestGetUser(t *testing.T) {
	testCases := []struct {
		name     string
		user     *coreModel.UserInfo
		excepted *model.GetUserResponse
	}{
		{
			name: "test_user_id",
			user: &coreModel.UserInfo{
				UserID:      "test_user_id",
				DisplayName: "test_display_name",
				Email:       "test_email",
				PhoneNumber: "test_phone_number",
				PhotoURL:    "test_photo_url",
			},
			excepted: &model.GetUserResponse{
				UserID:      "test_user_id",
				DisplayName: "test_display_name",
				Email:       "test_email",
				PhoneNumber: "test_phone_number",
				PhotoURL:    "test_photo_url",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			userController := NewUserController()
			httpRecorder := utils.RecordHandlerHttpRequest(
				userController.GetUser,
				"GET", "/"+testCase.name,
				nil,
				map[string]any{
					"user": testCase.user,
				},
			)

			if testCase.excepted != nil {
				assert.Equal(t, 200, httpRecorder.Code)
				assert.Equal(t, utils.ToJSONString(testCase.excepted), httpRecorder.Body.String())
			}
		})
	}
}
