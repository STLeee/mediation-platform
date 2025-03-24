package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// User is a user
type User struct {
	UserID      string    `json:"user_id" bson:"-"`
	FirebaseUID string    `json:"firebase_uid" bson:"firebase_uid"`
	DisplayName string    `json:"display_name" bson:"display_name"`
	Email       string    `json:"email" bson:"email"`
	PhoneNumber string    `json:"phone_number" bson:"phone_number"`
	PhotoURL    string    `json:"photo_url" bson:"photo_url"`
	Disabled    bool      `json:"disabled" bson:"disabled"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
	LastLoginAt time.Time `json:"last_login_at" bson:"last_login_at"`
}

// UserInMongoDB is a user in MongoDB
type UserInMongoDB struct {
	ID   bson.ObjectID `bson:"_id"`
	User `bson:",inline"`
}

func NewUserInMongoDB(user *User) (*UserInMongoDB, error) {
	var objectID bson.ObjectID
	var err error
	if user.UserID != "" {
		objectID, err = bson.ObjectIDFromHex(user.UserID)
		if err != nil {
			return nil, err
		}
	} else {
		objectID = bson.NewObjectID()
	}
	return &UserInMongoDB{
		ID:   objectID,
		User: *user,
	}, nil
}

func (userInMongoDB *UserInMongoDB) SetupDataFromDocument() error {
	userInMongoDB.User.UserID = userInMongoDB.ID.Hex()
	return nil
}
