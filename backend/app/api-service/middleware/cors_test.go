package middleware

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/STLeee/mediation-platform/backend/core/utils"
)

func TestCorsHandler(t *testing.T) {
	testCases := []struct {
		method        string
		excepted_code int
	}{
		{
			method:        "GET",
			excepted_code: 200,
		},
		{
			method:        "POST",
			excepted_code: 200,
		},
		{
			method:        "PUT",
			excepted_code: 200,
		},
		{
			method:        "DELETE",
			excepted_code: 200,
		},
		{
			method:        "OPTIONS",
			excepted_code: 204,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.method, func(t *testing.T) {
			httpRecorder := utils.RegisterAndRecordHttpRequest(func(routeGroup *gin.RouterGroup) {
				routeGroup.Use(CorsHandler())
				routeGroup.Handle(testCase.method, "/test", func(c *gin.Context) {
					c.JSON(200, gin.H{})
				})
			}, testCase.method, "/test", nil)

			assert.Equal(t, testCase.excepted_code, httpRecorder.Code)

			assert.Equal(t, "*", httpRecorder.Header().Get("Access-Control-Allow-Origin"))
			assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", httpRecorder.Header().Get("Access-Control-Allow-Methods"))
			assert.Equal(t, "Content-Type, Authorization", httpRecorder.Header().Get("Access-Control-Allow-Headers"))
			assert.Equal(t, "Authorization", httpRecorder.Header().Get("Access-Control-Expose-Headers"))
		})
	}
}
