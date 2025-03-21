package utils

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// ToJSONString converts an interface to a JSON string without error
func ToJSONString(v any) string {
	bytes, _ := json.Marshal(v)
	return string(bytes)
}

func ToObjectID(v string) bson.ObjectID {
	oid, _ := bson.ObjectIDFromHex(v)
	return oid
}
