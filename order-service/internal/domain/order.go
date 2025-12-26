package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/user/go-microservices/pkg/valueobject"
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
	ID            int64             `json:"id"`
	UserID        int64             `json:"user_id"`
	ProductID     int64             `json:"product_id"`
	ProductName   string            `json:"product_name"` // Snapshot
	UnitPrice     valueobject.Money `json:"unit_price"`   // Snapshot
	Quantity      int               `json:"quantity"`
	TotalPrice    valueobject.Money `json:"total_price"`
	OrderStatus   OrderStatus       `json:"order_status"`
	PaymentStatus PaymentStatus     `json:"payment_status"`
	CreatedAt     time.Time         `json:"created_at"`
}

// NewOrder is a factory function for the Order aggregate
func NewOrder(userID int64, productID int64, productName string, unitPrice valueobject.Money, quantity int) (*Order, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user id")
	}
	if quantity <= 0 {
		return nil, fmt.Errorf("quantity must be greater than zero")
	}

	return &Order{
		UserID:        userID,
		ProductID:     productID,
		ProductName:   productName,
		UnitPrice:     unitPrice,
		Quantity:      quantity,
		TotalPrice:    unitPrice.Multiply(quantity),
		OrderStatus:   OrderPending,
		PaymentStatus: PaymentPending,
		CreatedAt:     time.Now(),
	}, nil
}

// Pay marks the order as paid
func (o *Order) Pay() error {
	if o.OrderStatus == OrderCancelled {
		return fmt.Errorf("cannot pay for a cancelled order")
	}
	if o.PaymentStatus == PaymentPaid {
		return fmt.Errorf("order is already paid")
	}
	o.PaymentStatus = PaymentPaid
	return nil
}

// Complete completes the order
func (o *Order) Complete() error {
	if o.PaymentStatus != PaymentPaid {
		return fmt.Errorf("cannot complete an unpaid order")
	}
	o.OrderStatus = OrderCompleted
	return nil
}

// Cancel cancels the order
func (o *Order) Cancel() error {
	if o.OrderStatus == OrderCompleted {
		return fmt.Errorf("cannot cancel a completed order")
	}
	o.OrderStatus = OrderCancelled
	return nil
}

//go:generate mockery --name OrderRepository
type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id int64) (*Order, error)
	GetAll(ctx context.Context) ([]*Order, error)
	UpdateStatus(ctx context.Context, id int64, status OrderStatus, paymentStatus PaymentStatus) error
}

//go:generate mockery --name ProductClient
type ProductClient interface {
	GetProduct(ctx context.Context, id int64) (*ProductView, error)
	ReserveStock(ctx context.Context, id int64, qty int) error
	ReleaseStock(ctx context.Context, id int64, qty int) error
}

type ProductView struct {
	ID    int64             `json:"id"`
	Name  string            `json:"name"`
	Price valueobject.Money `json:"price"`
}
