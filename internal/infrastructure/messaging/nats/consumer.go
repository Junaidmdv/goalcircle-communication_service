package natsclient

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

type StreamConfig struct {
	Name        string
	Description string
	Subjects    []string
	MaxAge      time.Duration
}

type ConsumerConfig struct {
	Name          string
	DeliverPolicy jetstream.DeliverPolicy
	AckPolicy     jetstream.AckPolicy
	MaxDeliver    int
}

type JetStreamManager struct {
	js jetstream.JetStream
}

func NewJetStreamManager(js jetstream.JetStream) *JetStreamManager {
	return &JetStreamManager{
		js: js,
	}
}

func (m *JetStreamManager) CreateStream(
	ctx context.Context,
	config StreamConfig,
) (jetstream.Stream, error) {

	stream, err := m.js.CreateOrUpdateStream(
		ctx,
		jetstream.StreamConfig{
			Name:        config.Name,
			Description: config.Description,
			Subjects:    config.Subjects,
			Storage:     jetstream.FileStorage,
			Retention:   jetstream.LimitsPolicy,
			MaxAge:      config.MaxAge,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("create stream %s: %w", config.Name, err)
	}

	return stream, nil
}

func (m *JetStreamManager) CreateConsumer(
	ctx context.Context,
	streamName string,
	config ConsumerConfig,
) (jetstream.Consumer, error) {

	consumer, err := m.js.CreateOrUpdateConsumer(
		ctx,
		streamName,
		jetstream.ConsumerConfig{
			Durable:       config.Name,
			DeliverPolicy: config.DeliverPolicy,
			AckPolicy:     config.AckPolicy,
			MaxDeliver:    config.MaxDeliver,
		},
	)

	if err != nil {
		return nil, fmt.Errorf(
			"create consumer %s: %w",
			config.Name,
			err,
		)
	}

	return consumer, nil
}

func (m *JetStreamManager) Consume(ctx context.Context, consumer jetstream.Consumer) (jetstream.ConsumeContext, error) {

	cc, err := consumer.Consume(
		func(msg jetstream.Msg) {

		},
		jetstream.ConsumeErrHandler(func(consumeCtx jetstream.ConsumeContext, err error) {

		}))

	if err != nil {
		return nil, fmt.Errorf("failed consume message")
	}

	return cc, nil

}
