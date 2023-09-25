package repository

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
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

type repository struct {
	pool *pgxpool.Pool
}

func New() (Repository, error) {
	connString := "postgresql://postgres:password@localhost/orders"
	pool, err := pgxpool.New(context.TODO(), connString)
	if err != nil {
		return nil, err
	}

	return &repository{
		pool: pool,
	}, nil
}

func (r *repository) CreateOrder(ctx context.Context, order *models.Order) error {
	err := r.pool.QueryRow(ctx, `INSERT INTO orders (item_id, status, is_paid, notification_sent) VALUES ($1, 'PROCESSING', FALSE, FALSE) RETURNING id`, order.ItemID).Scan(&order.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateOrderStatus(ctx context.Context, orderID int, status string) error {
	return nil
}

func (r *repository) OrderSent(ctx context.Context, orderID int) error {
	_, err := r.pool.Exec(ctx, `UPDATE orders SET notification_sent = true, status = 'SUCCESS' WHERE id = $1`, orderID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) CancelOrder(ctx context.Context, orderID int) error {
	_, err := r.pool.Exec(ctx, `UPDATE orders SET status = 'FAIL' WHERE id = $1`, orderID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) OrderPaid(ctx context.Context, orderID int) error {
	_, err := r.pool.Exec(ctx, `UPDATE orders SET is_paid = TRUE WHERE id = $1`, orderID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetProcessingOrders(ctx context.Context) ([]*models.Order, error) {
	var mm []*models.Order

	err := pgxscan.Select(ctx, r.pool, &mm, `SELECT id, item_id, status, is_paid, notification_sent FROM orders WHERE status = 'PROCESSING'`)
	if err != nil {
		return nil, err
	}

	return mm, nil
}
