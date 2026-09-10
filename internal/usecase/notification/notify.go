package notification_usecase

import (
	"context"
	"time"

	"github.com/Junaidmdv/goalcircle-communication_service/internal/domain/entity"
	notificationrepo "github.com/Junaidmdv/goalcircle-communication_service/internal/domain/repository/notification"
	"github.com/Junaidmdv/goalcircle-communication_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-communication_service/pkg/logger"
	"github.com/google/uuid"
)

type NotificationUsecase interface {
	AddNotification(context.Context, *NotifyEvent) error
	GetNotifications(context.Context, *GetNotificationsReq) (*GetNotificationsRes, error)
	MarkNotificationAsRead(ctx context.Context, req MarkNotificationsAsReadReq) (*MarkNotificationAsReadRes, error)
	MarkAllNotificationAsRead(ctx context.Context, req MarkAllNotificationsAsReadReq) (*MarkAllNotificationsAsReadRes, error)
}

type notificationUsecase struct {
	notificationRepo notificationrepo.NotificationRepository
	logger           logger.Logger
}

func NewNotificationUsecase(nrepo notificationrepo.NotificationRepository, logger logger.Logger) NotificationUsecase {
	return &notificationUsecase{
		notificationRepo: nrepo,
		logger:           logger,
	}
}

func (u *notificationUsecase) AddNotification(ctx context.Context, input *NotifyEvent) error {

	if err := u.notificationRepo.AddNotification(ctx, &entity.Notification{
		ID:        uuid.New(),
		UserID:    input.UserID,
		Type:      input.Type,
		Title:     input.Title,
		Message:   input.Message,
		Data:      input.Data,
		CreatedAt: time.Now(),
	}); err != nil {
		return err
	}
	return nil
}

func (u *notificationUsecase) GetNotifications(ctx context.Context, input *GetNotificationsReq) (*GetNotificationsRes, error) {

	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		u.logger.Error("failed get user id", err)
		return nil, apperror.NewInternalError("something went wrong.Please try again later", err)
	}

	notification, err := u.notificationRepo.GetNotifications(ctx, userID)

	var notificationDetails []NotifficationDetails

	for _, n := range notification {
		notificationDetails = append(notificationDetails, NotifficationDetails{
			ID:        n.ID.String(),
			Type:      n.Type,
			Title:     n.Title,
			Message:   n.Message,
			Data:      n.Data,
			CreatedAt: n.CreatedAt,
		})
	}

	return &GetNotificationsRes{
		UserID:       notification[0].UserID.String(),
		Notification: notificationDetails,
	}, nil

}

func (u *notificationUsecase) MarkNotificationAsRead(ctx context.Context, req MarkNotificationsAsReadReq) (*MarkNotificationAsReadRes, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		u.logger.Error("failed get user id", err)
		return nil, apperror.NewInternalError("something went wrong.Please try again later", err)
	}

	notificationID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid notification id")
	}

	if err := u.notificationRepo.UpdatNotificationRead(ctx, notificationID, userID); err != nil {
		return nil, err
	}

	return &MarkNotificationAsReadRes{
		Success: true,
	}, nil

}

func (u *notificationUsecase) MarkAllNotificationAsRead(ctx context.Context, req MarkAllNotificationsAsReadReq) (*MarkAllNotificationsAsReadRes, error) {

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		u.logger.Error("failed get user id", err)
		return nil, apperror.NewInternalError("something went wrong.Please try again later", err)
	}

	count, err := u.notificationRepo.UpdatNotificationReadAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &MarkAllNotificationsAsReadRes{
		UpdatedCount: count,
	}, nil
}
