package model

type User struct {
	UserID      string `json:"user_id" bson:"-"`
	FirebaseUID string `json:"firebase_uid" bson:"firebase_uid"`
	DisplayName string `json:"display_name" bson:"display_name"`
	Email       string `json:"email" bson:"email"`
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
	PhotoURL    string `json:"photo_url" bson:"photo_url"`
	Disabled    bool   `json:"disabled" bson:"disabled"`
	CreatedAt   int64  `json:"created_at" bson:"created_at"`
	LastLoginAt int64  `json:"last_login_at" bson:"last_login_at"`
}
