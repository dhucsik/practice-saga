package subscriber

import (
	"payment-service/service"
	"payment-service/subscriber/nats"
)

type Subscriber interface {
	Subscribe() error
	Start() error
}

func New(srv service.Service) (Subscriber, error) {
	return nats.New(srv)
}
