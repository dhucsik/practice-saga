package repository

import (
	"context"
	"order-service/models"
)

type Repository interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	UpdateOrderStatus(ctx context.Context, orderID int, status string) error
	OrderPaid(ctx context.Context, orderID int) error
	CancelOrder(ctx context.Context, orderID int) error
	OrderSent(ctx context.Context, orderID int) error
	GetProcessingOrders(ctx context.Context) ([]*models.Order, error)
}

type repository struct{}

func New() Repository {
	return &repository{}
}

func (r *repository) CreateOrder(ctx context.Context, order *models.Order) error {
	order.ID = 10
	return nil
}

func (r *repository) UpdateOrderStatus(ctx context.Context, orderID int, status string) error {
	return nil
}

func (r *repository) OrderSent(ctx context.Context, orderID int) error {
	//order.NotificationSent = true
	return nil
}

func (r *repository) CancelOrder(ctx context.Context, orderID int) error {
	// order.Status = "Failed"
	return nil
}

func (r *repository) OrderPaid(ctx context.Context, orderID int) error {
	//order.IsPaid = true
	return nil
}

func (r *repository) GetProcessingOrders(ctx context.Context) ([]*models.Order, error) {
	return nil, nil
}
