package natsclient

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

type Publisher struct {
	js jetstream.JetStream
}

func NewPublisher(js jetstream.JetStream) *Publisher {
	return &Publisher{
		js: js,
	}
}

func (p *Publisher) Publish(
	ctx context.Context,
	subject string,
	event any,
) (*jetstream.PubAck, error) {

	data, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("marshal event: %w", err)
	}

	ack, err := p.js.Publish(ctx, subject, data)
	if err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}

	return ack, nil
}