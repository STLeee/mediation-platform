package router

import (
	"testing"

	"github.com/STLeee/mediation-platform/backend/core/utils"
)

func TestRegisterHealthRouter(t *testing.T) {
	utils.TestRouterRegister(t, RegisterHealthRouter, []string{
		"/liveness",
		"/readiness",
	})
}
