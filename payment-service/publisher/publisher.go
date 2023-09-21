package publisher

import (
	"context"
	"payment-service/publisher/nats"
)

type Publisher interface {
	Publish(ctx context.Context, data []byte, subject string) error
}

func New() (Publisher, error) {
	return nats.New()
}
