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

	// Decode the token
	decodedToken, err := DecodeMockFirebaseIDToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, decodedToken)
	assert.Equal(t, projectID, decodedToken["aud"])
	assert.Equal(t, uid, decodedToken["uid"])
}

func TestDecodeMockFirebaseIDToken_InvalidToken(t *testing.T) {
	token := "invalid"
	decodedToken, err := DecodeMockFirebaseIDToken(token)
	assert.Error(t, err)
	assert.Nil(t, decodedToken)
}
