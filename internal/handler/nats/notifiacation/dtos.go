package notification_natshandler

import "github.com/google/uuid"

type NotifyEvent struct {
	UserID  uuid.UUID      `json:"user_id"`
	Type    string         `json:"type"`
	Title   string         `json:"title"`
	Message string         `json:"message"`
	Data    map[string]any `json:"data,omitempty"`
}


