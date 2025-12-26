package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/user/go-microservices/order-service/internal/domain"
	pkgerrors "github.com/user/go-microservices/pkg/errors"
	"github.com/user/go-microservices/pkg/logger"
	"go.uber.org/zap"
)

type postgresRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) domain.OrderRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, o *domain.Order) error {
	query := `
		INSERT INTO orders (user_id, product_id, product_name, unit_price, quantity, total_price, order_status, payment_status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	now := time.Now().UTC()
	err := r.db.QueryRowContext(ctx, query,
		o.UserID, o.ProductID, o.ProductName, o.UnitPrice, o.Quantity, o.TotalPrice,
		o.OrderStatus, o.PaymentStatus, now,
	).Scan(&o.ID)

	if err != nil {
		logger.FromContext(ctx).Error("failed to create order", zap.Error(err))
		return pkgerrors.ErrInternal
	}
	o.CreatedAt = now
	return nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (*domain.Order, error) {
	query := `SELECT id, user_id, product_id, product_name, unit_price, quantity, total_price, order_status, payment_status, created_at FROM orders WHERE id = $1`

	o := &domain.Order{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&o.ID, &o.UserID, &o.ProductID, &o.ProductName, &o.UnitPrice,
		&o.Quantity, &o.TotalPrice, &o.OrderStatus, &o.PaymentStatus, &o.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, pkgerrors.ErrNotFound
	}
	if err != nil {
		logger.FromContext(ctx).Error("failed to get order", zap.Error(err))
		return nil, pkgerrors.ErrInternal
	}
	return o, nil
}

func (r *postgresRepository) UpdateStatus(ctx context.Context, id int64, status domain.OrderStatus, paymentStatus domain.PaymentStatus) error {
	query := `UPDATE orders SET order_status = $1, payment_status = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, status, paymentStatus, id)
	if err != nil {
		logger.FromContext(ctx).Error("failed to update order status", zap.Error(err))
		return pkgerrors.ErrInternal
	}
	return nil
}

func (r *postgresRepository) GetAll(ctx context.Context) ([]*domain.Order, error) {
	query := `SELECT id, user_id, product_id, product_name, unit_price, quantity, total_price, order_status, payment_status, created_at FROM orders`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		logger.FromContext(ctx).Error("failed to get all orders", zap.Error(err))
		return nil, pkgerrors.ErrInternal
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		o := &domain.Order{}
		err := rows.Scan(
			&o.ID, &o.UserID, &o.ProductID, &o.ProductName, &o.UnitPrice,
			&o.Quantity, &o.TotalPrice, &o.OrderStatus, &o.PaymentStatus, &o.CreatedAt,
		)
		if err != nil {
			logger.FromContext(ctx).Error("failed to scan order", zap.Error(err))
			return nil, pkgerrors.ErrInternal
		}
		orders = append(orders, o)
	}
	return orders, nil
}
