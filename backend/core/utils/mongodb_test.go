package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertStringToObjectID(t *testing.T) {
	testCases := []struct {
		name     string
		id       string
		expected any
	}{
		{
			name:     "valid-id",
			id:       "5f4b8f1f9d1e4b0001f3f3b1",
			expected: "5f4b8f1f9d1e4b0001f3f3b1",
		},
		{
			name:     "empty-id",
			id:       "",
			expected: "000000000000000000000000",
		},
		{
			name:     "invalid-id",
			id:       "invalid-id",
			expected: "000000000000000000000000",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := ConvertStringToObjectID(testCase.id); got != testCase.expected {
				assert.Equal(t, testCase.expected, got.Hex())
			}
		})
	}
}
