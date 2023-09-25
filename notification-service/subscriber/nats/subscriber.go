package nats

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"notification-service/models"
	"notification-service/service"
)

type broker struct {
	nc  *nats.Conn
	srv service.Service
}

func New(srv service.Service) (*broker, error) {
	nc, err := nats.Connect("")
	if err != nil {
		return nil, err
	}

	return &broker{
		nc:  nc,
		srv: srv,
	}, nil
}

func (b *broker) SubscribeOrders() error {
	b.nc.Subscribe("notifications.orders", func(msg *nats.Msg) {
		var notifiaction models.Notification

		if err := json.Unmarshal(msg.Data, &notifiaction); err != nil {
			log.Println("error: " + err.Error())
			return
		}

		err := b.srv.SendNotification(context.TODO(), &notifiaction)
		if err != nil {
			log.Println("error: " + err.Error())
			return
		}
	})

	return nil
}

func (b *broker) Start() error {
	if err := b.SubscribeOrders(); err != nil {
		return err
	}

	return nil
}
