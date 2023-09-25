package service

import (
	"context"
	"encoding/json"
	"payment-service/models"
	"payment-service/publisher"
	"payment-service/repository"
)

type Service interface {
	ProcessPayment(ctx context.Context, order *models.Order) error
	PaidPayment(ctx context.Context, paymentID int) error
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

func (s *service) PushToOrders(ctx context.Context, payment *models.Payment) error {
	data, err := json.Marshal(payment)
	if err != nil {
		return err
	}

	err = s.pub.Publish(ctx, data, "orders.payments")
	if err != nil {
		return err
	}
	return nil
}

func (s *service) ProcessPayment(ctx context.Context, order *models.Order) error {
	//make a payment
	payment, err := s.repo.CheckOrder(ctx, order.ID)
	if err != nil {
		return err
	}

	if payment.OrderID != 0 {
		return s.PushToOrders(ctx, payment)
	}

	payment = &models.Payment{OrderID: order.ID}
	err = s.repo.CreatePayment(ctx, payment)
	if err != nil {
		payment.Status = "FAIL"
		return err
	}
	// if success publish to order service payment information with status success
	payment.Status = "SUCCESS"
	return s.PushToOrders(ctx, payment)
}

func (s *service) PaidPayment(ctx context.Context, paymentID int) error {
	return s.repo.PaidPayment(ctx, paymentID)
}
