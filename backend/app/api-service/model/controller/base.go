package controller

// MessageResponse is a response for message
type MessageResponse struct {
	Message string `json:"message"`
}

// NewMessageResponse creates a new MessageResponse
func NewMessageResponse(message string) *MessageResponse {
	return &MessageResponse{Message: message}
}
