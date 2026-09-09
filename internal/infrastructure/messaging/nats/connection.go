package natsclient

import (
	"time"

	"github.com/Junaidmdv/goalcircle-communication_service/internal/config"
	"github.com/Junaidmdv/goalcircle-communication_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-communication_service/pkg/logger"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Connection struct {
	conn   *nats.Conn
	js     jetstream.JetStream
	logger logger.Logger
}

func NewConnection(cfg config.NatsConfig, logger logger.Logger) (*Connection, error) {
	conn, err := nats.Connect(
		cfg.URL,
		nats.Name("communication-service"),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			logger.Error("NATS disconnected:", "error", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			logger.Info("NATS reconnected", "time", time.November.String())
		}),
	)

	if err != nil {
		logger.Error("failed connect nats", "error", err)
		return nil, apperror.NewInternalError("something went wrong.Please try again later", err)
	}

	return &Connection{
		conn: conn,
	}, nil
}

func (c *Connection) Close() {
	if c.conn != nil {
		c.conn.Drain()
	}
}
