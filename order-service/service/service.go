package service

import (
	"context"
	"encoding/json"
	"order-service/models"
	"order-service/publisher"
	"order-service/repository"
)

type Service interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	PushToPayment(ctx context.Context, order *models.Order) error
	UpdateOrderPayment(ctx context.Context, payment *models.Payment) error
	PushToNotification(ctx context.Context, order *models.Order) error
	CancelOrder(ctx context.Context, orderID int) error
	GetProcessingOrders(ctx context.Context) ([]*models.Order, error)
	UpdateOrderNotification(ctx context.Context, notification *models.Notification) error
}

type service struct {
	repo repository.Repository
	pub  publisher.Publisher
}

func New(repo repository.Repository, pub publisher.Publisher) Service {
	return &service{
		repo: repo,
		pub:  pub,
	}
}

func (s *service) CreateOrder(ctx context.Context, order *models.Order) error {
	// Order {status: processing}
	order.Status = "Processing"
	err := s.repo.CreateOrder(ctx, order)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) PushToPayment(ctx context.Context, order *models.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	err = s.pub.Publish(ctx, data, "payments.orders")
	if err != nil {
		return s.CancelOrder(ctx, order.ID)
	}

	return nil
}

func (s *service) UpdateOrderPayment(ctx context.Context, payment *models.Payment) error {
	if payment.Status == "SUCCESS" {
		return s.repo.OrderPaid(ctx, payment.OrderID)
	}

	if payment.Status == "FAIL" {
		return s.CancelOrder(ctx, payment.OrderID)
	}

	return nil
}

func (s *service) PushToNotification(ctx context.Context, order *models.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	err = s.pub.Publish(ctx, data, "notifications.orders")
	if err != nil {
		return err
	}

	return nil
}

func (s *service) UpdateOrderNotification(ctx context.Context, notification *models.Notification) error {
	return s.repo.OrderSent(ctx, notification.ID)
}

func (s *service) GetProcessingOrders(ctx context.Context) ([]*models.Order, error) {
	return s.repo.GetProcessingOrders(ctx)
}

//
//func (s *service) UpdateOrder(ctx context.Context, payment *models.Payment) error {
//	// if payment success then => status = success
//	if payment.Status == "success" {
//		err := s.repo.UpdateOrderStatus(ctx, payment.OrderID, "success")
//		if err != nil {
//			return err
//		}
//
//		return nil
//	}
//	err := s.repo.UpdateOrderStatus(ctx, payment.OrderID, "failed")
//	if err != nil {
//		return err
//	}
//
//	return nil
//
//	// else status = failed
//
//}

func (s *service) CancelOrder(ctx context.Context, orderID int) error {
	return s.repo.CancelOrder(ctx, orderID)
}

//func (s *service) PerformOrder(ctx context.Context, order *models.Order) error {
//	if !order.IsCreated {
//		err := s.CreateOrder(ctx, order)
//		if err != nil {
//			return err
//		}
//	}
//
//	if !order.IsPaid {
//		err := s.PushToPayment(ctx, order)
//		if err != nil {
//			if err := s.CancelOrder(ctx, order.ID); err != nil {
//				log.Printf("error: %v", err.Error())
//			}
//
//			return err
//		}
//	}
//
//	if order.IsCreated && order.IsPaid {
//		err := s.PushToNotification(ctx, order)
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
