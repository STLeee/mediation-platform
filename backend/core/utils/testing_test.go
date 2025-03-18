package utils

import (
	"bytes"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTestRouterRegister(t *testing.T) {
	routerRegisterFunc := func(router *gin.RouterGroup) {
		router.GET("/a", func(ctx *gin.Context) {})
		router.GET("/b", func(ctx *gin.Context) {})
	}
	TestRouterRegister(t, routerRegisterFunc, []string{
		"/a",
		"/b",
	})
}

func TestRegisterAndRecordHttpRequest(t *testing.T) {
	routerRegisterFunc := func(router *gin.RouterGroup) {
		router.GET("/test", func(ctx *gin.Context) {
			body := make([]byte, ctx.Request.ContentLength)
			ctx.Request.Body.Read(body)
			ctx.JSON(200, gin.H{"message": string(body)})
		})
	}
	body := bytes.NewReader([]byte("this is testing"))
	httpRecorder := RegisterAndRecordHttpRequest(routerRegisterFunc, "GET", "/test", body)
	assert.NotNil(t, httpRecorder)
	assert.Equal(t, 200, httpRecorder.Code)
	assert.Equal(t, "{\"message\":\"this is testing\"}", httpRecorder.Body.String())
}
