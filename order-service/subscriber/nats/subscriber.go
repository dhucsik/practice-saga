package nats

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"order-service/models"
	"order-service/service"
)

type broker struct {
	nc  *nats.Conn
	srv service.Service
}

func New(srv service.Service) (*broker, error) {
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		return nil, err
	}

	return &broker{nc: nc, srv: srv}, nil
}

func (b *broker) SubscribePayments() error {
	b.nc.Subscribe("orders.payments", func(msg *nats.Msg) {
		var payment models.Payment
		err := json.Unmarshal(msg.Data, &payment)
		if err != nil {
			log.Println("error: " + err.Error())
			return
		}

		err = b.srv.UpdateOrderPayment(context.TODO(), &payment)
		if err != nil {
			log.Println("error: " + err.Error())
		}
	})

	return nil
}

func (b *broker) SubscribeNotifications() error {
	b.nc.Subscribe("orders.notifications", func(msg *nats.Msg) {
		var notification models.Notification
		err := json.Unmarshal(msg.Data, &notification)
		if err != nil {
			log.Println("error: " + err.Error())
			return
		}

		err = b.srv.UpdateOrderNotification(context.TODO(), &notification)
		if err != nil {
			log.Println("error: " + err.Error())
		}
	})

	return nil
}

func (b *broker) Start() error {
	if err := b.SubscribePayments(); err != nil {
		return err
	}

	if err := b.SubscribeNotifications(); err != nil {
		return err
	}

	return nil
}
