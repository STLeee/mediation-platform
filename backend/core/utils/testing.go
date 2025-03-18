package utils

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestRouterRegister tests the router registration
func TestRouterRegister(t *testing.T, routerRegisterFunc func(*gin.RouterGroup), exceptedPaths []string) {
	engine := gin.Default()
	routerRegisterFunc(engine.Group("/"))

	routes := engine.Routes()
	var paths []string
	for _, route := range routes {
		paths = append(paths, route.Path)
	}

	assert.ElementsMatch(t, exceptedPaths, paths)
}

// RegisterAndRecordHttpRequest registers the router and records the http request
func RegisterAndRecordHttpRequest(routerRegisterFunc func(*gin.RouterGroup), method string, path string, body io.Reader) *httptest.ResponseRecorder {
	engine := gin.Default()
	routerRegisterFunc(engine.Group("/"))

	return RecordHttpRequest(engine, method, path, body)
}

// RecordHttpRequest records the http request
func RecordHttpRequest(engine *gin.Engine, method string, path string, body io.Reader) *httptest.ResponseRecorder {
	httpRecorder := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, body)
	engine.ServeHTTP(httpRecorder, req)

	return httpRecorder
}
