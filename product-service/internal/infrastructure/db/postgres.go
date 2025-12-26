package repository

import (
	"context"
	"database/sql"
	"time"

	pkgerrors "github.com/user/go-microservices/pkg/errors"
	"github.com/user/go-microservices/pkg/logger"
	"github.com/user/go-microservices/product-service/internal/domain"
	"go.uber.org/zap"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) domain.ProductRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, p *domain.Product) error {
	query := `
		INSERT INTO products (sku, name, description, price, total_qty, reserved_qty, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, 0, $6, $7, $8)
		RETURNING id`

	now := time.Now().UTC()
	err := r.db.QueryRowContext(ctx, query, p.SKU, p.Name, p.Description, p.Price, p.TotalQty, p.IsActive, now, now).Scan(&p.ID)
	if err != nil {
		logger.FromContext(ctx).Error("failed to create product", zap.Error(err))
		return pkgerrors.ErrInternal
	}
	p.CreatedAt = now
	p.UpdatedAt = now
	return nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	query := `SELECT id, sku, name, description, price, total_qty, reserved_qty, is_active, created_at, updated_at FROM products WHERE id = $1`

	p := &domain.Product{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.SKU, &p.Name, &p.Description, &p.Price,
		&p.TotalQty, &p.ReservedQty, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, pkgerrors.ErrNotFound
	}
	if err != nil {
		logger.FromContext(ctx).Error("failed to get product", zap.Error(err))
		return nil, pkgerrors.ErrInternal
	}
	return p, nil
}

func (r *postgresRepository) ReserveStock(ctx context.Context, id int64, qty int) error {
	query := `
		UPDATE products 
		SET reserved_qty = reserved_qty + $1, updated_at = NOW()
		WHERE id = $2 AND (total_qty - reserved_qty) >= $1
	`
	res, err := r.db.ExecContext(ctx, query, qty, id)
	if err != nil {
		logger.FromContext(ctx).Error("failed to reserve stock", zap.Error(err))
		return pkgerrors.ErrInternal
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return pkgerrors.ErrInternal
	}
	if rows == 0 {
		return pkgerrors.ErrInsufficientStock
	}
	return nil
}

func (r *postgresRepository) ReleaseStock(ctx context.Context, id int64, qty int) error {
	query := `
		UPDATE products 
		SET reserved_qty = reserved_qty - $1, updated_at = NOW()
		WHERE id = $2 AND reserved_qty >= $1
	`
	res, err := r.db.ExecContext(ctx, query, qty, id)
	if err != nil {
		logger.FromContext(ctx).Error("failed to release stock", zap.Error(err))
		return pkgerrors.ErrInternal
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		// Should not happen if logic is correct, but safety check
		return pkgerrors.ErrInternal
	}
	return nil
}

func (r *postgresRepository) ConfirmStock(ctx context.Context, id int64, qty int) error {
	// Confirm means we permanently remove from global stock and reduce reserved
	query := `
		UPDATE products 
		SET total_qty = total_qty - $1, reserved_qty = reserved_qty - $1, updated_at = NOW()
		WHERE id = $2 AND reserved_qty >= $1
	`
	res, err := r.db.ExecContext(ctx, query, qty, id)
	if err != nil {
		logger.FromContext(ctx).Error("failed to confirm stock", zap.Error(err))
		return pkgerrors.ErrInternal
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return pkgerrors.ErrInternal
	}
	return nil
}

func (r *postgresRepository) GetAll(ctx context.Context) ([]*domain.Product, error) {
	query := `SELECT id, sku, name, description, price, total_qty, reserved_qty, is_active, created_at, updated_at FROM products`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		logger.FromContext(ctx).Error("failed to get all products", zap.Error(err))
		return nil, pkgerrors.ErrInternal
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		p := &domain.Product{}
		err := rows.Scan(
			&p.ID, &p.SKU, &p.Name, &p.Description, &p.Price,
			&p.TotalQty, &p.ReservedQty, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			logger.FromContext(ctx).Error("failed to scan product", zap.Error(err))
			return nil, pkgerrors.ErrInternal
		}
		products = append(products, p)
	}
	return products, nil
}
