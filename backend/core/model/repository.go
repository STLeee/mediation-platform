package model

const (
	UPDATED_TIMESTAMP_FIELD = "updated_at"
)

// MongoDBData interface for data in MongoDB
type MongoDBDocument interface {
	SetupDataFromDocument() error
}
