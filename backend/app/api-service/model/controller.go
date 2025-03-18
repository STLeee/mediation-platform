package model

// MessageResponse is a response for message
type MessageResponse struct {
	Message string `json:"message" example:"ok"`
}

// NewMessageResponse creates a new MessageResponse
func NewMessageResponse(message string) *MessageResponse {
	return &MessageResponse{Message: message}
}
