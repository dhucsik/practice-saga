package repository

import (
	"context"
	"payment-service/models"
)

type Repository interface {
	CreatePayment(ctx context.Context, payment *models.Payment) error
	CheckOrder(ctx context.Context, orderID int) (*models.Payment, error)
}

type repository struct{}

func New() Repository {
	return &repository{}
}

func (r *repository) CreatePayment(ctx context.Context, payment *models.Payment) error {
	payment.ID = 13
	return nil
}

func (r *repository) CheckOrder(ctx context.Context, orderID int) (*models.Payment, error) {
	// check if order already paid or not
	return nil, nil
}
