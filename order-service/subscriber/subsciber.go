package subscriber

import (
	"order-service/service"
	"order-service/subscriber/nats"
)

type Subscriber interface {
	SubscribePayments() error
	SubscribeNotifications() error
	Start() error
}

func New(srv service.Service) (Subscriber, error) {
	return nats.New(srv)
}
