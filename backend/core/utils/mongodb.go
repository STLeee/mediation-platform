package utils

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ConvertStringToObjectID converts a string to an MongoDB ObjectID
func ConvertStringToObjectID(v string) bson.ObjectID {
	oid, _ := bson.ObjectIDFromHex(v)
	return oid
}
