package notification_natshandler

import (
	"context"
	"encoding/json"

	notification_usecase "github.com/Junaidmdv/goalcircle-communication_service/internal/usecase/notification"
	"github.com/Junaidmdv/goalcircle-communication_service/pkg/logger"
	"github.com/nats-io/nats.go/jetstream"
)

type NotificationNatsHandler interface {
	Handle(msg jetstream.Msg)
}

type notificationNatsHandler struct {
	notificationUsecase notification_usecase.NotificationUsecase
	logger              logger.Logger
}

func NewNotificationNatsHandler() NotificationNatsHandler {
	return &notificationNatsHandler{}
}

func (h *notificationNatsHandler) Handle(msg jetstream.Msg) {

	ctx := context.Background()

	var evt NotifyEvent
	if err := json.Unmarshal(msg.Data(), &evt); err != nil {
		h.logger.Error("invalid notify payload", "error", err)
		_ = msg.Term()
		return
	}

	event, err := ToNotifyEvent(&evt)
	if err != nil {
		_ = msg.Term()
		return
	}

	if err := h.notificationUsecase.AddNotification(ctx, event); err != nil {
		h.logger.Error("failed to process notify event", "error", err)
		_ = msg.Nak()
		return
	}
}
