package entity

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationTypePlayerRequest   NotificationType = "player_request"
	NotificationTypeRequestAccepted NotificationType = "request_accepted"
	NotificationTypeRequestRejected NotificationType = "request_rejected"
	NotificationTypeTeamInvitation  NotificationType = "team_invitation"
)

type Notification struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	UserID    uuid.UUID
	Type      NotificationType
	Title     string
	Message   string
	Data      map[string]any `gorm:"serializer:json;type:jsonb"`
	IsRead    bool
	ReadAt    *time.Time
	CreatedAt time.Time
}
