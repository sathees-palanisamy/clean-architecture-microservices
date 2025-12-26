package domain

import (
	"context"
	"time"
)

type OrderStatus string
type PaymentStatus string

const (
	OrderPending   OrderStatus = "PENDING"
	OrderCompleted OrderStatus = "COMPLETED"
	OrderCancelled OrderStatus = "CANCELLED"

	PaymentPending PaymentStatus = "PENDING"
	PaymentPaid    PaymentStatus = "PAID"
	PaymentFailed  PaymentStatus = "FAILED"
)

type Order struct {
	ID            int64         `json:"id"`
	UserID        int64         `json:"user_id"`
	ProductID     int64         `json:"product_id"`
	ProductName   string        `json:"product_name"` // Snapshot
	UnitPrice     float64       `json:"unit_price"`   // Snapshot
	Quantity      int           `json:"quantity"`
	TotalPrice    float64       `json:"total_price"`
	OrderStatus   OrderStatus   `json:"order_status"`
	PaymentStatus PaymentStatus `json:"payment_status"`
	CreatedAt     time.Time     `json:"created_at"`
}

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id int64) (*Order, error)
	GetAll(ctx context.Context) ([]*Order, error)
	UpdateStatus(ctx context.Context, id int64, status OrderStatus, paymentStatus PaymentStatus) error
}

type ProductClient interface {
	GetProduct(ctx context.Context, id int64) (*ProductView, error)
	ReserveStock(ctx context.Context, id int64, qty int) error
	ReleaseStock(ctx context.Context, id int64, qty int) error
}

type ProductView struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
