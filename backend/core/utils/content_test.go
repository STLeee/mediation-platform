package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToJSONString(t *testing.T) {
	type args struct {
		v interface{}
	}
	testCases := []struct {
		name string
		args args
		want string
	}{
		{
			name: "map",
			args: args{
				v: map[string]interface{}{
					"key": "value",
				},
			},
			want: "{\"key\":\"value\"}",
		},
		{
			name: "slice",
			args: args{
				v: []string{"value1", "value2"},
			},
			want: "[\"value1\",\"value2\"]",
		},
		{
			name: "string",
			args: args{
				v: "value",
			},
			want: "\"value\"",
		},
		{
			name: "int",
			args: args{
				v: 1,
			},
			want: "1",
		},
		{
			name: "float",
			args: args{
				v: 1.1,
			},
			want: "1.1",
		},
		{
			name: "bool",
			args: args{
				v: true,
			},
			want: "true",
		},
		{
			name: "nil",
			args: args{
				v: nil,
			},
			want: "null",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := ConvertToJSONString(testCase.args.v); got != testCase.want {
				assert.Equal(t, testCase.want, got)
			}
		})
	}
}

func TestConvertStructToMap(t *testing.T) {
	testCases := []struct {
		name       string
		v          interface{}
		expected   map[string]interface{}
		errMessage string
	}{
		{
			name: "valid-input/struct",
			v: struct {
				Key string
			}{
				Key: "value",
			},
			expected: map[string]interface{}{
				"Key": "value",
			},
		},
		{
			name: "valid-input/struct-pointer",
			v: &struct {
				Key string
			}{
				Key: "value",
			},
			expected: map[string]interface{}{
				"Key": "value",
			},
		},
		{
			name:       "invalid-input",
			v:          "invalid",
			errMessage: "not a struct or pointer to struct",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := ConvertStructToMap(testCase.v)
			if err != nil {
				assert.Equal(t, testCase.errMessage, err.Error())
			} else {
				assert.Equal(t, testCase.expected, got)
			}
		})
	}
}
