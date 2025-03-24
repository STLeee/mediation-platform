package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMockFirebaseIDToken(t *testing.T) {
	projectID := "test-project-id"
	uid := "test-uid"
	token := CreateMockFirebaseIDToken(projectID, uid)
	assert.NotNil(t, token)
}
