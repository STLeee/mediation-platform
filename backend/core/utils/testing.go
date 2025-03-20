package utils

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestEngineRouterRegister tests the engine router registration
func TestEngineRouterRegister(t *testing.T, routerRegisterFunc func(*gin.Engine), exceptedPaths []string) {
	engine := gin.Default()
	routerRegisterFunc(engine)

	routes := engine.Routes()
	var paths []string
	for _, route := range routes {
		paths = append(paths, route.Path)
	}

	assert.ElementsMatch(t, exceptedPaths, paths)
}

// TestRouterRegister tests the router registration
func TestRouterRegister(t *testing.T, routerRegisterFunc func(*gin.RouterGroup), exceptedPaths []string) {
	TestEngineRouterRegister(t, func(engine *gin.Engine) {
		routerRegisterFunc(engine.Group("/"))
	}, exceptedPaths)
}

// RegisterAndRecordHttpRequest registers the router and records the http request
func RegisterAndRecordHttpRequest(routerRegisterFunc func(*gin.RouterGroup), method string, path string, body io.Reader) *httptest.ResponseRecorder {
	httpRecorder := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(httpRecorder)
	routerRegisterFunc(engine.Group("/"))
	ctx.Request = httptest.NewRequest(method, path, body)

	engine.ServeHTTP(httpRecorder, ctx.Request)
	return httpRecorder
}

// RecordHandlerHttpRequest records the http request with handler
func RecordHandlerHttpRequest(handler gin.HandlerFunc, method string, path string, body io.Reader, ctxValues map[string]any) *httptest.ResponseRecorder {
	httpRecorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(httpRecorder)
	for key, value := range ctxValues {
		ctx.Set(key, value)
	}
	ctx.Request = httptest.NewRequest(method, path, body)

	handler(ctx)
	return httpRecorder
}
