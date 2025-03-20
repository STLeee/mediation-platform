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

func TestRegisterV1UserRouter(t *testing.T) {
	utils.TestRouterRegister(t, RegisterV1UserRouter, []string{
		"/:user_id",
	})
}
