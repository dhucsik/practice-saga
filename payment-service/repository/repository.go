package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"payment-service/models"
)

type Repository interface {
	CreatePayment(ctx context.Context, payment *models.Payment) error
	CheckOrder(ctx context.Context, orderID int) (*models.Payment, error)
	PaidPayment(ctx context.Context, paymentID int) error
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
	err := r.pool.QueryRow(ctx, `INSERT INTO payments (order_id, status) VALUES ($1, 'PENDING') RETURNING id`, payment.OrderID).Scan(&payment.ID)
	if err != nil {
		return fmt.Errorf("create payment", err)
	}

	return nil
}

func (r *repository) CheckOrder(ctx context.Context, orderID int) (*models.Payment, error) {
	payment := &models.Payment{}
	err := r.pool.QueryRow(ctx, `SELECT id, order_id, status FROM payments WHERE id = $1`, orderID).Scan(&payment.ID, &payment.OrderID, &payment.Status)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("check order", err)
	}

	return payment, nil
}

func (r *repository) PaidPayment(ctx context.Context, paymentID int) error {
	_, err := r.pool.Exec(ctx, `UPDATE payments SET status = 'SUCCESS' WHERE id = $1`, paymentID)
	if err != nil {
		return err
	}

	return nil
}
