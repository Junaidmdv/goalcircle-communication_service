package notification_usecase

import (
	"time"

	"github.com/Junaidmdv/goalcircle-communication_service/internal/domain/entity"
	"github.com/google/uuid"
)

type NotifyEvent struct {
	UserID  uuid.UUID
	Type    entity.NotificationType
	Title   string
	Message string
	Data    map[string]any
}

type GetNotificationsReq struct {
	UserID string
}

type NotifficationDetails struct {
	ID        string
	Type      entity.NotificationType
	Title     string
	Message   string
	Data      map[string]any
	CreatedAt time.Time
}

type GetNotificationsRes struct {
	UserID       string
	Notification []NotifficationDetails
}

type MarkNotificationsAsReadReq struct {
	ID     string
	UserID string
}

type MarkNotificationAsReadRes struct {
	Success bool
}

type MarkAllNotificationsAsReadReq struct {
	UserID string
}

type MarkAllNotificationsAsReadRes struct {
	UpdatedCount int64
}
