package domain

import (
	"context"
	"time"

	"github.com/user/go-microservices/pkg/valueobject"
)

type Product struct {
	ID          int64             `json:"id"`
	SKU         string            `json:"sku"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Price       valueobject.Money `json:"price"`
	TotalQty    int               `json:"total_qty"`
	ReservedQty int               `json:"reserved_qty"`
	IsActive    bool              `json:"is_active"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

func (p *Product) AvailableQty() int {
	return p.TotalQty - p.ReservedQty
}

//go:generate mockery --name ProductRepository
type ProductRepository interface {
	Create(ctx context.Context, p *Product) error
	GetByID(ctx context.Context, id int64) (*Product, error)
	ReserveStock(ctx context.Context, id int64, qty int) error
	ReleaseStock(ctx context.Context, id int64, qty int) error
	ConfirmStock(ctx context.Context, id int64, qty int) error
	GetAll(ctx context.Context) ([]*Product, error)
}
