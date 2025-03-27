package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/STLeee/mediation-platform/backend/core/db"
	"github.com/STLeee/mediation-platform/backend/core/model"
)

var LocalRepositoryConfigs = &RepositoryConfigs{
	UserDB: &MongoDBRepositoryConfig{
		Database:   "mediation-platform",
		Collection: "user",
	},
}

// RepositoryErrorType struct for repository error type
type RepositoryErrorType string

const (
	RepositoryErrorTypeServerError    RepositoryErrorType = "server_error"
	RepositoryErrorTypeConfigError    RepositoryErrorType = "config_error"
	RepositoryErrorTypeRecordNotFound RepositoryErrorType = "record_not_found"
	RepositoryErrorTypeInvalidID      RepositoryErrorType = "invalid_id"
	RepositoryErrorTypeInvalidData    RepositoryErrorType = "invalid_data"
)

var RepositoryErrorDefaultMessages = map[RepositoryErrorType]string{
	RepositoryErrorTypeServerError:    "server error",
	RepositoryErrorTypeConfigError:    "config error",
	RepositoryErrorTypeRecordNotFound: "record not found",
	RepositoryErrorTypeInvalidID:      "invalid ID",
	RepositoryErrorTypeInvalidData:    "invalid data",
}

// RepositoryError struct for repository error
type RepositoryError struct {
	ErrType    RepositoryErrorType
	Database   string
	Collection string
	Message    string
	Err        error
}

// Error returns the error message
func (e RepositoryError) Error() string {
	message := e.Message
	if message == "" {
		if defaultMessage, ok := RepositoryErrorDefaultMessages[e.ErrType]; ok {
			message = defaultMessage
		}
	}
	message = fmt.Sprintf("%s/%s: %s", e.Database, e.Collection, message)
	if e.Err != nil {
		message = strings.Join([]string{message, e.Err.Error()}, ": ")
	}
	return message
}

// Unwrap returns the wrapped error
func (e RepositoryError) Unwrap() error {
	return e.Err
}

// RepositoryName struct for repository name
type RepositoryName string

const (
	RepositoryNameUserDB RepositoryName = "user_db"
)

// MongoDBRepositoryConfigs struct for MongoDB repository configs
type RepositoryConfigs struct {
	UserDB *MongoDBRepositoryConfig `yaml:"user_db"`
}

// MongoDBRepositoryConfig struct for MongoDB repository config
type MongoDBRepositoryConfig struct {
	Database   string `yaml:"database"`
	Collection string `yaml:"collection"`
}

// MongoDBRepository struct for MongoDB repository
type MongoDBRepository struct {
	mongoDB    *db.MongoDB
	collection *mongo.Collection
	cfg        *MongoDBRepositoryConfig
}

// NewMongoDBRepository creates a new MongoDB repository
func NewMongoDBRepository(mongoDB *db.MongoDB, cfg *MongoDBRepositoryConfig) *MongoDBRepository {
	return &MongoDBRepository{
		mongoDB:    mongoDB,
		collection: mongoDB.Database(cfg.Database).Collection(cfg.Collection),
		cfg:        cfg,
	}
}

// InsertOne inserts one
func (repo *MongoDBRepository) InsertOne(ctx context.Context, data model.MongoDBDocument) (string, error) {
	res, err := repo.collection.InsertOne(ctx, data)
	if err != nil {
		return "", RepositoryError{
			ErrType:    RepositoryErrorTypeServerError,
			Database:   repo.cfg.Database,
			Collection: repo.cfg.Collection,
			Message:    "failed to insert one",
			Err:        err,
		}
	}
	return res.InsertedID.(bson.ObjectID).Hex(), nil
}

// FindOneByID finds one by ID
func (repo *MongoDBRepository) FindByID(ctx context.Context, id string, result model.MongoDBDocument) error {
	// Convert ID to ObjectID
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return RepositoryError{
			ErrType:    RepositoryErrorTypeInvalidID,
			Database:   repo.cfg.Database,
			Collection: repo.cfg.Collection,
			Message:    "invalid ID",
		}
	}

	// Find one
	filter := bson.M{"_id": objectID}
	return repo.FindOneByFilter(ctx, filter, result)
}

// FindOneByFilter finds one by filter
func (repo *MongoDBRepository) FindOneByFilter(ctx context.Context, filter map[string]any, result model.MongoDBDocument) error {
	// Find one
	if err := repo.collection.FindOne(ctx, filter).Decode(result); err != nil {
		if err == mongo.ErrNoDocuments {
			return RepositoryError{
				ErrType:    RepositoryErrorTypeRecordNotFound,
				Database:   repo.cfg.Database,
				Collection: repo.cfg.Collection,
				Message:    "record not found",
			}
		}
		return RepositoryError{
			ErrType:    RepositoryErrorTypeServerError,
			Database:   repo.cfg.Database,
			Collection: repo.cfg.Collection,
			Message:    "failed to get one by ID",
			Err:        err,
		}
	}

	// Setup data from document
	err := result.SetupDataFromDocument()
	if err != nil {
		return RepositoryError{
			ErrType:    RepositoryErrorTypeInvalidData,
			Database:   repo.cfg.Database,
			Collection: repo.cfg.Collection,
			Message:    "setup data from document failed",
			Err:        err,
		}
	}
	return nil
}

func (repo *MongoDBRepository) UpdateByID(ctx context.Context, id string, data map[string]any) error {
	// Convert ID to ObjectID
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return RepositoryError{
			ErrType:    RepositoryErrorTypeInvalidID,
			Database:   repo.cfg.Database,
			Collection: repo.cfg.Collection,
			Message:    "invalid ID",
		}
	}

	// Set update data
	data[model.UpdatedTimestampFieldName] = time.Now()
	update := bson.M{"$set": data}

	// Update one
	res, err := repo.collection.UpdateByID(ctx, objectID, update)
	if err != nil {
		return RepositoryError{
			ErrType:    RepositoryErrorTypeServerError,
			Database:   repo.cfg.Database,
			Collection: repo.cfg.Collection,
			Message:    "failed to update one by ID",
			Err:        err,
		}
	}
	if res.MatchedCount == 0 {
		return RepositoryError{
			ErrType:    RepositoryErrorTypeRecordNotFound,
			Database:   repo.cfg.Database,
			Collection: repo.cfg.Collection,
			Message:    "record not found",
		}
	}
	return nil
}

func (repo *MongoDBRepository) DeleteByID(ctx context.Context, id string) error {
	// Convert ID to ObjectID
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return RepositoryError{
			ErrType:    RepositoryErrorTypeInvalidID,
			Database:   repo.cfg.Database,
			Collection: repo.cfg.Collection,
			Message:    "invalid ID",
		}
	}

	// Delete one
	filter := bson.M{"_id": objectID}
	res, err := repo.collection.DeleteOne(ctx, filter)
	if err != nil {
		return RepositoryError{
			ErrType:    RepositoryErrorTypeServerError,
			Database:   repo.cfg.Database,
			Collection: repo.cfg.Collection,
			Message:    "failed to delete one by ID",
			Err:        err,
		}
	}
	if res.DeletedCount == 0 {
		return RepositoryError{
			ErrType:    RepositoryErrorTypeRecordNotFound,
			Database:   repo.cfg.Database,
			Collection: repo.cfg.Collection,
			Message:    "record not found",
		}
	}
	return nil
}
