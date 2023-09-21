package service

import (
	"context"
	"encoding/json"
	"mail-service/models"
	"mail-service/publisher"
)

type Service interface {
	SendNotification(ctx context.Context, notification *models.Notification) error
}

type service struct {
	pub publisher.Publisher
}

func New(pub publisher.Publisher) Service {
	return &service{
		pub: pub,
	}
}

func (s *service) PushToOrders(ctx context.Context, notification *models.Notification) error {
	data, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	err = s.pub.Publish(ctx, data, "orders.notifications")
	if err != nil {
		return err
	}

	return nil
}

func (s *service) SendNotification(ctx context.Context, notification *models.Notification) error {
	// send notification

	return s.PushToOrders(ctx, notification)
}
