package notification_repo

import (
	"context"
	"time"

	"github.com/Junaidmdv/goalcircle-communication_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-communication_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-communication_service/pkg/logger"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	AddNotification(context.Context, *entity.Notification) error
	GetNotifications(ctx context.Context, userID uuid.UUID) ([]entity.Notification, error)
	UpdatNotificationRead(ctx context.Context, ID uuid.UUID, UserID uuid.UUID) error
	UpdatNotificationReadAll(ctx context.Context, userID uuid.UUID) (int64, error)
}

type notificationRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewNotificationRepository(db *gorm.DB, logger logger.Logger) NotificationRepository {
	return &notificationRepository{
		db:     db,
		logger: logger,
	}
}

func (n *notificationRepository) AddNotification(ctx context.Context, data *entity.Notification) error {

	if err := n.db.WithContext(ctx).Create(data).Error; err != nil {
		n.logger.Error("database error", "error", err)
		return apperror.NewInternalError("something went wrong. Please try again later", err)
	}
	return nil
}

func (n *notificationRepository) GetNotifications(ctx context.Context, userID uuid.UUID) ([]entity.Notification, error) {

	var notify []entity.Notification
	result := n.db.WithContext(ctx).Where("user_id=? AND is_read=?", userID, false).Order("created_at DESC").Find(&notify)

	if result.Error != nil {
		return nil, apperror.NewInternalError("internal server error", result.Error)
	}

	return notify, nil
}

func (n *notificationRepository) UpdatNotificationRead(ctx context.Context, ID uuid.UUID, UserID uuid.UUID) error {

	result := n.db.WithContext(ctx).Model(&entity.Notification{}).Where("id=? AND user_id=? AND is_read=", ID, UserID, false).Updates(map[string]interface{}{
		"is_read": true,
		"read_at": time.Now(),
	})

	if result.Error != nil {
		n.logger.Error("databse error", "error", result.Error, "method", "notificationRepo.UpdatNotificationRead")
		return apperror.NewInternalError("something went wrong.Please try again later", result.Error)
	}
	return nil
}

func (n *notificationRepository) UpdatNotificationReadAll(ctx context.Context, userID uuid.UUID) (int64, error) {

	result := n.db.WithContext(ctx).Model(&entity.Notification{}).Where("user_id=? AND is_read=?", userID, false).Updates(map[string]interface{}{
		"is_read": true,
		"read_at": time.Now(),
	})

	if result.Error != nil {
		n.logger.Error("databse error", "error", result.Error, "method", "notificationRepo.UpdatNotificationReadAll")
		return -1, apperror.NewInternalError("something went wrong.Please try again later", result.Error)
	}
	return result.RowsAffected, nil
}
