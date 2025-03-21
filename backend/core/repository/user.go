package repository

import (
	"context"

	"github.com/STLeee/mediation-platform/backend/core/db"
	"github.com/STLeee/mediation-platform/backend/core/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// UserInMongoDB is a user in MongoDB
type UserInMongoDB struct {
	ID         bson.ObjectID `bson:"_id"`
	model.User `bson:",inline"`
}

func NewUserInMongoDB(user *model.User) (*UserInMongoDB, error) {
	userInMongoDB := &UserInMongoDB{
		User: *user,
	}
	if user.UserID != "" {
		var err error
		userInMongoDB.ID, err = bson.ObjectIDFromHex(user.UserID)
		if err != nil {
			return nil, err
		}
	}
	return userInMongoDB, nil
}

func (userInMongoDB *UserInMongoDB) ToUser() *model.User {
	userInMongoDB.UserID = userInMongoDB.ID.Hex()
	return &userInMongoDB.User
}

// UserRepository is an interface for user repository
type UserRepository interface {
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
}

// UserMongoDBRepository is a MongoDB repository for user
type UserMongoDBRepository struct {
	MongoDBRepository
}

// NewUserMongoDB creates a new UserMongoDB
func NewUserMongoDB(mongoDB *db.MongoDB, cfg *MongoDBRepositoryConfig) *UserMongoDBRepository {
	return &UserMongoDBRepository{
		MongoDBRepository: *NewMongoDBRepository(mongoDB, cfg),
	}
}

// GetUserByID get a user by user ID
func (repo *UserMongoDBRepository) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	userInMongoDB := &UserInMongoDB{}
	err := repo.FindByID(ctx, userID, userInMongoDB)
	if err != nil {
		return nil, err
	}
	return userInMongoDB.ToUser(), nil
}
