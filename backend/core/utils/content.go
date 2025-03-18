package utils

import "encoding/json"

// ToJSONString converts an interface to a JSON string without error
func ToJSONString(v interface{}) string {
	bytes, _ := json.Marshal(v)
	return string(bytes)
}
