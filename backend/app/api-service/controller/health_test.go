package controller

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	controllerModel "github.com/STLeee/mediation-platform/backend/app/api-service/model/controller"
	"github.com/STLeee/mediation-platform/backend/core/utils"
)

func TestHealthControllerLiveness(t *testing.T) {
	routerRegisterFunc := func(r *gin.RouterGroup) {
		healthController := NewHealthController()
		r.GET("/liveness", healthController.Liveness)
	}
	httpRecorder := utils.RegisterAndRecordHttpRequest(routerRegisterFunc, "GET", "/liveness", nil)

	assert.Equal(t, 200, httpRecorder.Code)
	expectedResponse := controllerModel.NewMessageResponse("ok")
	assert.Equal(t, utils.ToJSONString(expectedResponse), httpRecorder.Body.String())
}
func TestHealthControllerReadiness(t *testing.T) {
	routerRegisterFunc := func(r *gin.RouterGroup) {
		healthController := NewHealthController()
		r.GET("/readiness", healthController.Readiness)
	}
	httpRecorder := utils.RegisterAndRecordHttpRequest(routerRegisterFunc, "GET", "/readiness", nil)

	assert.Equal(t, 200, httpRecorder.Code)
	expectedResponse := controllerModel.NewMessageResponse("ok")
	assert.Equal(t, utils.ToJSONString(expectedResponse), httpRecorder.Body.String())
}
