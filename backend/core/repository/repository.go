package repository

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/STLeee/mediation-platform/backend/core/db"
)

var LocalMongoDBRepositoryConfigs = MongoDBRepositoryConfigs{
	RepositoryNameUser: {
		Database:   "mediation-platform",
		Collection: "user",
	},
}

// RepositoryErrorType struct for repository error type
type RepositoryErrorType string

const (
	RepositoryErrorTypeServerError    RepositoryErrorType = "server_error"
	RepositoryErrorTypeRecordNotFound RepositoryErrorType = "record_not_found"
	RepositoryErrorTypeInvalidID      RepositoryErrorType = "invalid_id"
)

var RepositoryErrorDefaultMessages = map[RepositoryErrorType]string{
	RepositoryErrorTypeServerError:    "server error",
	RepositoryErrorTypeRecordNotFound: "record not found",
	RepositoryErrorTypeInvalidID:      "invalid ID",
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
	RepositoryNameUser RepositoryName = "user"
)

// MongoDBRepositoryConfigs struct for MongoDB repository configs
type MongoDBRepositoryConfigs map[RepositoryName]*MongoDBRepositoryConfig

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
		collection: mongoDB.Collection(cfg.Database, cfg.Collection),
		cfg:        cfg,
	}
}

// FindOneByID finds one by ID
func (repo *MongoDBRepository) FindByID(ctx context.Context, id string, result any) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return RepositoryError{
			ErrType:    RepositoryErrorTypeInvalidID,
			Database:   repo.cfg.Database,
			Collection: repo.cfg.Collection,
			Message:    "invalid ID",
		}
	}
	filter := bson.M{"_id": objectID}
	if err := repo.collection.FindOne(ctx, filter).Decode(result); err != nil {
		fmt.Printf("err: %+v\n", err)
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
	fmt.Printf("result: %+v\n", result)
	return nil
}
