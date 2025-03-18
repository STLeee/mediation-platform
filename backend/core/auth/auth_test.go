package auth

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthError(t *testing.T) {
	ae := AuthError{
		authName: "test",
		err:      fmt.Errorf("test error"),
	}
	assert.Equal(t, "[auth - test] test error", ae.Error())
	assert.Equal(t, "test error", ae.Unwrap().Error())
}
