package notification_natshandler

import (
	"fmt"

	"github.com/Junaidmdv/goalcircle-communication_service/internal/domain/entity"
	nuc "github.com/Junaidmdv/goalcircle-communication_service/internal/usecase/notification"
)

func MapNotificationType(value string) (entity.NotificationType, error) {

	switch value {
	case string(entity.NotificationTypePlayerRequest):
		return entity.NotificationTypePlayerRequest, nil

	case string(entity.NotificationTypeRequestAccepted):
		return entity.NotificationTypeRequestAccepted, nil

	case string(entity.NotificationTypeRequestRejected):
		return entity.NotificationTypeRequestRejected, nil

	case string(entity.NotificationTypeTeamInvitation):
		return entity.NotificationTypeTeamInvitation, nil

	default:
		return "", fmt.Errorf("unsupported notification type: %s", value)
	}
}

func ToNotifyEvent(req *NotifyEvent) (*nuc.NotifyEvent, error) {

	notifytype, err := MapNotificationType(req.Type)
	if err != nil {
		return nil, err
	}

	return &nuc.NotifyEvent{
		UserID:  req.UserID,
		Type:    notifytype,
		Title:   req.Title,
		Message: req.Message,
	}, nil
}
