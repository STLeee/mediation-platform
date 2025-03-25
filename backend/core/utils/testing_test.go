package utils

import (
	"bytes"
	"testing"
	"time"

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

func TestSimplyValidTimestamp(t *testing.T) {
	testCases := []struct {
		name      string
		timestamp time.Time
		expected  bool
	}{
		{
			name:      "valid-timestamp",
			timestamp: time.Now(),
			expected:  true,
		},
		{
			name:      "zero-timestamp",
			timestamp: time.Time{},
			expected:  false,
		},
		{
			name:      "future-timestamp",
			timestamp: time.Now().Add(time.Hour),
			expected:  false,
		},
		{
			name:      "past-timestamp",
			timestamp: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected:  false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			isValid := SimplyValidTimestamp(testCase.timestamp)
			assert.Equal(t, testCase.expected, isValid)
		})
	}
}
