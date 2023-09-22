package nats

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"payment-service/models"
	"payment-service/service"
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

	return &broker{nc: nc, srv: srv}, nil
}

func (b *broker) SubscribeOrders() error {
	// order data
	b.nc.Subscribe("payments.orders", func(msg *nats.Msg) {
		var order *models.Order

		if err := json.Unmarshal(msg.Data, order); err != nil {
			log.Println("error: " + err.Error())
			return
		}

		err := b.srv.ProcessPayment(context.TODO(), order)
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
