package utils

import "encoding/json"

func ToJSONString(v interface{}) string {
	bytes, _ := json.Marshal(v)
	return string(bytes)
}
