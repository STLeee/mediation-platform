package controller

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMessageResponse(t *testing.T) {
	message := "test"
	got := NewMessageResponse(message)
	assert.Equal(t, message, got.Message)

	jsonString := `{"message":"test"}`
	bytes, _ := json.Marshal(got)
	assert.Equal(t, jsonString, string(bytes))
}
