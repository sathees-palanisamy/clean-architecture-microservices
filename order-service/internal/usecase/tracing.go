package usecase

import (
	"context"

	"github.com/user/go-microservices/order-service/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type tracingOrderUsecase struct {
	next   OrderUsecase
	tracer trace.Tracer
}

func NewTracingOrderUsecase(next OrderUsecase) OrderUsecase {
	return &tracingOrderUsecase{
		next:   next,
		tracer: otel.Tracer("order-usecase"),
	}
}

func (u *tracingOrderUsecase) CreateOrder(ctx context.Context, userID, productID int64, qty int) (*domain.Order, error) {
	ctx, span := u.tracer.Start(ctx, "CreateOrder")
	defer span.End()
	return u.next.CreateOrder(ctx, userID, productID, qty)
}

func (u *tracingOrderUsecase) GetOrder(ctx context.Context, id int64) (*domain.Order, error) {
	ctx, span := u.tracer.Start(ctx, "GetOrder")
	defer span.End()
	return u.next.GetOrder(ctx, id)
}

func (u *tracingOrderUsecase) GetAllOrders(ctx context.Context) ([]*domain.Order, error) {
	ctx, span := u.tracer.Start(ctx, "GetAllOrders")
	defer span.End()
	return u.next.GetAllOrders(ctx)
}

func (u *tracingOrderUsecase) CancelOrder(ctx context.Context, id int64) error {
	ctx, span := u.tracer.Start(ctx, "CancelOrder")
	defer span.End()
	return u.next.CancelOrder(ctx, id)
}
