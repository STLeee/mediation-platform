package middleware

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/STLeee/mediation-platform/backend/app/api-service/model"
	"github.com/STLeee/mediation-platform/backend/core/utils"
)

func TestErrorHandler(t *testing.T) {
	testCases := []struct {
		name          string
		err           error
		excepted_code int
	}{
		{
			name:          "no-error",
			err:           nil,
			excepted_code: http.StatusOK,
		},
		{
			name:          "bad-request",
			err:           model.HttpStatusCodeError{StatusCode: http.StatusBadRequest},
			excepted_code: http.StatusBadRequest,
		},
		{
			name:          "internal-server-error",
			err:           model.HttpStatusCodeError{StatusCode: http.StatusInternalServerError},
			excepted_code: http.StatusInternalServerError,
		},
		{
			name:          "unknown-error",
			err:           fmt.Errorf("test-error"),
			excepted_code: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			httpRecorder := utils.RegisterAndRecordHttpRequest(func(routeGroup *gin.RouterGroup) {
				routeGroup.Use(ErrorHandler())
				routeGroup.Handle("GET", "/test", func(c *gin.Context) {
					if testCase.err != nil {
						c.Error(testCase.err)
						return
					}
					c.JSON(http.StatusOK, gin.H{})
				})
			}, "GET", "/test", nil)

			assert.Equal(t, testCase.excepted_code, httpRecorder.Code)
		})
	}
}
