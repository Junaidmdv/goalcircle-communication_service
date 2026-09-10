package notification_handler

import (
	"context"
	"time"

	nuc "github.com/Junaidmdv/goalcircle-communication_service/internal/usecase/notification"
	"github.com/Junaidmdv/goalcircle-communication_service/pkg/apperror"
	compb "github.com/Junaidmdv/goalcircle-protos/communication/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

type NotificationHandler struct {
	notificationUsecase nuc.NotificationUsecase
	timeOut             time.Duration
	compb.UnimplementedNotificationServiceServer
}

func NewNotificationHandler(notificationusecase nuc.NotificationUsecase, timeOut time.Duration) *NotificationHandler {
	return &NotificationHandler{
		notificationUsecase: notificationusecase,
		timeOut:             timeOut,
	}
}

func (nh *NotificationHandler) GetNotifications(ctx context.Context, input *compb.GetNotificationsReq) (*compb.GetNotificationsRes, error) {

	ctx, cancel := context.WithTimeout(ctx, nh.timeOut)
	defer cancel()

	res, err := nh.notificationUsecase.GetNotifications(ctx, &nuc.GetNotificationsReq{
		UserID: input.UserId,
	})

	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	var listNotification []compb.GetNotificationsRes_NotificationDetail

	for _, r := range res.Notification {
		data, err := structpb.NewStruct(r.Data)
		if err != nil {
			return nil, err
		}
		listNotification = append(listNotification, compb.GetNotificationsRes_NotificationDetail{
			NotificationId:   r.ID,
			NotificationType: string(r.Type),
			Title:            r.Title,
			Message:          r.Message,
			Data:             data,
		})
	}

	return &compb.GetNotificationsRes{
		UserId: res.UserID,
	}, nil
}

func (nh *NotificationHandler) MarkAllNotificationsAsRead(ctx context.Context, input *compb.MarkAllNotificationsAsReadRequest) (*compb.MarkAllNotificationsAsReadResponse, error) {

	res, err := nh.notificationUsecase.MarkAllNotificationAsRead(ctx, nuc.MarkAllNotificationsAsReadReq{
		UserID: input.UserId,
	})
	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	return &compb.MarkAllNotificationsAsReadResponse{
		UpdateCount: res.UpdatedCount,
	}, nil
}

func (nh *NotificationHandler) MarkNotificationAsRead(ctx context.Context, input *compb.MarkNotificationAsReadRequest) (*compb.MarkNotificationAsReadResponse, error) {
	res, err := nh.notificationUsecase.MarkNotificationAsRead(ctx, nuc.MarkNotificationsAsReadReq{
		ID:     input.NotificationId,
		UserID: input.UserId,
	})
	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	return &compb.MarkNotificationAsReadResponse{
		Success: res.Success,
	}, nil
}

func (nh *NotificationHandler) AcceptPlayerRequest(ctx context.Context, input *compb.AcceptPlayerRequestRequest) (*compb.AcceptPlayerRequestResponse, error) {
	return nil, nil
}

func (nh *NotificationHandler) RejectPlayerRequest(ctx context.Context, input *compb.RejectPlayerRequestRequest) (*compb.RejectPlayerRequestResponse, error) {
	return nil, nil
}
