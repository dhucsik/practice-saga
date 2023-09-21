package worker

import (
	"context"
	"github.com/robfig/cron/v3"
	"log"
	"order-service/service"
)

type PusherWorker interface {
}

type pusherWorker struct {
	srv service.Service
}

func New(srv service.Service) PusherWorker {
	return &pusherWorker{
		srv: srv,
	}
}

func (w *pusherWorker) Start(ctx context.Context) {
	c := cron.New()

	_, err := c.AddFunc("@every 10s", func() {
		wErr := w.Work()
		if wErr != nil {
			log.Println("error: ", wErr.Error())
		}
	})
	if err != nil {
		log.Println("error: ", err.Error())
	}

	go c.Start()
}

func (w *pusherWorker) Work() error {
	ctx := context.TODO()
	orders, err := w.srv.GetProcessingOrders(ctx)
	if err != nil {
		return err
	}

	for _, order := range orders {
		if !order.IsPaid {
			err := w.srv.PushToPayment(ctx, order)
			if err != nil {
				log.Println("error: ", err.Error())
			}
		}

		if order.IsPaid && !order.NotificationSent {
			err := w.srv.PushToNotification(ctx, order)
			if err != nil {
				log.Println("error: ", err.Error())
			}
		}
	}

	return nil
}
