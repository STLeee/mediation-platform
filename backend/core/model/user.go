package model

type UserInfo struct {
	UserID        string `json:"user_id"`
	FirebaseUID   string `json:"firebase_uid"`
	DisplayName   string `json:"display_name"`
	Email         string `json:"email"`
	PhoneNumber   string `json:"phone_number"`
	PhotoURL      string `json:"photo_url"`
	Disabled      bool   `json:"disabled"`
	EmailVerified bool   `json:"email_verified"`
}
