package nats

import (
	"context"
	"github.com/nats-io/nats.go"
)

type broker struct {
	nc *nats.Conn
}

func New() (*broker, error) {
	nc, err := nats.Connect("")
	if err != nil {
		return nil, err
	}

	return &broker{nc: nc}, nil
}

func (b *broker) Publish(ctx context.Context, data []byte, subject string) error {
	return b.nc.Publish(subject, data)
}
