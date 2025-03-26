package repository

import (
	"context"
	"time"

	"github.com/STLeee/mediation-platform/backend/core/auth"
	"github.com/STLeee/mediation-platform/backend/core/db"
	"github.com/STLeee/mediation-platform/backend/core/model"
)

// UserRepository is an interface for user repository
type UserRepository interface {
	GetUserByToken(ctx context.Context, token string) (*model.User, error)
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
}

// UserMongoDBRepository is a MongoDB repository for user
type UserMongoDBRepository struct {
	MongoDBRepository
	authService auth.BaseAuthService
}

// NewUserMongoDBRepository creates a new UserMongoDBRepository
func NewUserMongoDBRepository(mongoDB *db.MongoDB, cfg *MongoDBRepositoryConfig) *UserMongoDBRepository {
	return &UserMongoDBRepository{
		MongoDBRepository: *NewMongoDBRepository(mongoDB, cfg),
	}
}

// SetAuthService sets an auth service
func (repo *UserMongoDBRepository) SetAuthService(authService auth.BaseAuthService) {
	repo.authService = authService
}

// GetUserByToken get a user by token
func (repo *UserMongoDBRepository) GetUserByToken(ctx context.Context, token string) (*model.User, error) {
	if repo.authService == nil {
		return nil, RepositoryError{
			ErrType: RepositoryErrorTypeConfigError,
			Message: "auth service is not set",
		}
	}

	// Get user ID from auth service
	authUID, err := repo.authService.AuthenticateByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Get user data from auth service
	userFromAuth, mapping, err := repo.authService.GetUserInfoAndMapping(ctx, authUID)
	if err != nil {
		return nil, err
	}

	// Get user from MongoDB
	user, err := repo.GetUserByFilter(ctx, mapping)
	if err != nil {
		// If user not found, create a new user
		if _, ok := err.(RepositoryError); ok {
			userID, err := repo.CreateUser(ctx, userFromAuth)
			if err != nil {
				return nil, err
			}
			user, err = repo.GetUserByID(ctx, userID)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return user, nil
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

func (repo *UserMongoDBRepository) GetUserByFilter(ctx context.Context, filter map[string]any) (*model.User, error) {
	// Find by filter
	userInMongoDB := &model.UserInMongoDB{}
	err := repo.FindOneByFilter(ctx, filter, userInMongoDB)
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
