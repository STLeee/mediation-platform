package model

const (
	UpdatedTimestampFieldName = "updated_at"
)

// MongoDBData interface for data in MongoDB
type MongoDBDocument interface {
	SetupDataFromDocument() error
}
