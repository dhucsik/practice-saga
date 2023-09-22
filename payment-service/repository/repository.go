package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"payment-service/models"
)

type Repository interface {
	CreatePayment(ctx context.Context, payment *models.Payment) error
	CheckOrder(ctx context.Context, orderID int) (*models.Payment, error)
}

type repository struct {
	pool *pgxpool.Pool
}

func New() (Repository, error) {
	connString := "postgresql://postgres:password@localhost/payments"
	pool, err := pgxpool.New(context.TODO(), connString)
	if err != nil {
		return nil, err
	}

	return &repository{
		pool: pool,
	}, nil
}

func (r *repository) CreatePayment(ctx context.Context, payment *models.Payment) error {
	err := r.pool.QueryRow(ctx, "INSERT INTO payments (order_id) VALUES ($1) RETURNING id").Scan(&payment.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) CheckOrder(ctx context.Context, orderID int) (*models.Payment, error) {
	payment := &models.Payment{}
	err := r.pool.QueryRow(ctx, "SELECT id, order_id FROM payments WHERE order_id = $1", orderID).Scan(payment)
	if err != nil {
		return nil, err
	}

	return payment, nil
}
