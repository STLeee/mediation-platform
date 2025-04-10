package repository

import (
	"context"
	"time"

	"github.com/STLeee/mediation-platform/backend/core/auth"
	"github.com/STLeee/mediation-platform/backend/core/db"
	"github.com/STLeee/mediation-platform/backend/core/model"
)

// UserDBRepository is an interface for user repository in database
type UserDBRepository interface {
	CreateUser(ctx context.Context, user *model.User) (string, error)
	GetUserByAuthUID(ctx context.Context, authName auth.AuthServiceName, authUID string) (*model.User, error)
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
}

// UserMongoDBRepository is a MongoDB repository for user
type UserMongoDBRepository struct {
	MongoDBRepository
}

// NewUserMongoDBRepository creates a new UserMongoDBRepository
func NewUserMongoDBRepository(mongoDB *db.MongoDB, cfg *MongoDBRepositoryConfig) *UserMongoDBRepository {
	return &UserMongoDBRepository{
		MongoDBRepository: *NewMongoDBRepository(mongoDB, cfg),
	}
}

// CreateUser create a user
func (repo *UserMongoDBRepository) CreateUser(ctx context.Context, user *model.User) (string, error) {
	// Set created at, updated at, and last login at
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.LastLoginAt = now

	// Insert one
	userInMongoDB, err := model.NewUserInMongoDB(user)
	if err != nil {
		return "", err
	}
	return repo.InsertOne(ctx, userInMongoDB)
}

// GetUserByAuthUID get a user by auth UID
func (repo *UserMongoDBRepository) GetUserByAuthUID(ctx context.Context, authName auth.AuthServiceName, authUID string) (*model.User, error) {
	authUIDFilter := map[string]any{}
	switch authName {
	case auth.AuthServiceNameFirebase:
		authUIDFilter["firebase_uid"] = authUID
	default:
		return nil, RepositoryError{
			ErrType: RepositoryErrorTypeServerError,
			Message: "unsupported auth service",
		}
	}

	// Find by filter
	userInMongoDB := &model.UserInMongoDB{}
	err := repo.FindOneByFilter(ctx, authUIDFilter, userInMongoDB)
	if err != nil {
		return nil, err
	}
	return &userInMongoDB.User, nil
}

// GetUserByID get a user by user ID
func (repo *UserMongoDBRepository) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	// Find by ID
	userInMongoDB := &model.UserInMongoDB{}
	err := repo.FindByID(ctx, userID, userInMongoDB)
	if err != nil {
		return nil, err
	}
	return &userInMongoDB.User, nil
}

// UpdateUser updates a user
func (repo *UserMongoDBRepository) UpdateUserByID(ctx context.Context, userID string, updateData map[string]any) error {
	return repo.UpdateByID(ctx, userID, updateData)
}

// DeleteUserByID deletes a user by user ID
func (repo *UserMongoDBRepository) DeleteUserByID(ctx context.Context, userID string) error {
	// Delete by ID
	return repo.DeleteByID(ctx, userID)
}
