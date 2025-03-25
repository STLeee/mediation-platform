package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// ConvertToJSONString converts an interface to a JSON string without error
func ConvertToJSONString(v any) string {
	bytes, _ := json.Marshal(v)
	return string(bytes)
}

// ConvertStructToMap converts a struct to a map
func ConvertStructToMap(in interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return out, fmt.Errorf("not a struct or pointer to struct")
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		out[t.Field(i).Name] = v.Field(i).Interface()
	}
	return out, nil
}
