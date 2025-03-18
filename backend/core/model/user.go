package model

type UserInfo struct {
	UID           string `json:"uid"`
	DisplayName   string `json:"display_name"`
	Email         string `json:"email"`
	PhoneNumber   string `json:"phone_number"`
	PhotoURL      string `json:"photo_url"`
	Disabled      bool   `json:"disabled"`
	EmailVerified bool   `json:"email_verified"`
}
