package model

// MessageResponse is a response for message
type MessageResponse struct {
	Message string `json:"message" example:"ok"`
}

type GetUserResponse struct {
	UserID      string `json:"user_id" example:"1234567890"`
	DisplayName string `json:"display_name" example:"Scott Li"`
	Email       string `json:"email" example:"example@mediation-platform.com"`
	PhoneNumber string `json:"phone_number" example:"+886987654321"`
	PhotoURL    string `json:"photo_url" example:"https://example.com/photo.jpg"`
}
