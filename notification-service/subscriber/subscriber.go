package subscriber

import (
	"mail-service/service"
	"mail-service/subscriber/nats"
)

type Subscriber interface {
	SubscribeOrders() error
	Start() error
}

func New(srv service.Service) (Subscriber, error) {
	return nats.New(srv)
}
