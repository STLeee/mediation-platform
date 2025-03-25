package utils

import (
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestEngineRouterRegister tests the engine router registration
func TestEngineRouterRegister(t *testing.T, routerRegisterFunc func(*gin.Engine), expectedPaths []string) {
	gin.SetMode(gin.TestMode)
	engine := gin.Default()
	routerRegisterFunc(engine)

	routes := engine.Routes()
	var paths []string
	for _, route := range routes {
		paths = append(paths, route.Path)
	}

	assert.ElementsMatch(t, expectedPaths, paths)
}

// TestRouterRegister tests the router registration
func TestRouterRegister(t *testing.T, routerRegisterFunc func(*gin.RouterGroup), expectedPaths []string) {
	TestEngineRouterRegister(t, func(engine *gin.Engine) {
		routerRegisterFunc(engine.Group("/"))
	}, expectedPaths)
}

// RegisterAndRecordHttpRequest registers the router and records the http request
func RegisterAndRecordHttpRequest(routerRegisterFunc func(*gin.RouterGroup), method string, path string, body io.Reader) *httptest.ResponseRecorder {
	httpRecorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, engine := gin.CreateTestContext(httpRecorder)
	routerRegisterFunc(engine.Group("/"))
	ctx.Request = httptest.NewRequest(method, path, body)

	engine.ServeHTTP(httpRecorder, ctx.Request)
	return httpRecorder
}

// SimplyValidTimestamp simply checks if the timestamp is within a valid range
func SimplyValidTimestamp(timestamp time.Time) bool {
	fromTime := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)
	toTime := time.Now()
	return (fromTime.Before(timestamp) || fromTime.Equal(timestamp)) && (toTime.After(timestamp) || toTime.Equal(timestamp))
}
