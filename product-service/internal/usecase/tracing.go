package usecase

import (
	"context"

	"github.com/user/go-microservices/product-service/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type tracingProductUsecase struct {
	next   ProductUsecase
	tracer trace.Tracer
}

func NewTracingProductUsecase(next ProductUsecase) ProductUsecase {
	return &tracingProductUsecase{
		next:   next,
		tracer: otel.Tracer("product-usecase"),
	}
}

func (u *tracingProductUsecase) CreateProduct(ctx context.Context, p *domain.Product) error {
	ctx, span := u.tracer.Start(ctx, "CreateProduct")
	defer span.End()
	return u.next.CreateProduct(ctx, p)
}

func (u *tracingProductUsecase) GetProduct(ctx context.Context, id int64) (*domain.Product, error) {
	ctx, span := u.tracer.Start(ctx, "GetProduct")
	defer span.End()
	return u.next.GetProduct(ctx, id)
}

func (u *tracingProductUsecase) ReserveStock(ctx context.Context, id int64, qty int) error {
	ctx, span := u.tracer.Start(ctx, "ReserveStock")
	defer span.End()
	return u.next.ReserveStock(ctx, id, qty)
}

func (u *tracingProductUsecase) ReleaseStock(ctx context.Context, id int64, qty int) error {
	ctx, span := u.tracer.Start(ctx, "ReleaseStock")
	defer span.End()
	return u.next.ReleaseStock(ctx, id, qty)
}

func (u *tracingProductUsecase) ConfirmStock(ctx context.Context, id int64, qty int) error {
	ctx, span := u.tracer.Start(ctx, "ConfirmStock")
	defer span.End()
	return u.next.ConfirmStock(ctx, id, qty)
}

func (u *tracingProductUsecase) GetAllProducts(ctx context.Context) ([]*domain.Product, error) {
	ctx, span := u.tracer.Start(ctx, "GetAllProducts")
	defer span.End()
	return u.next.GetAllProducts(ctx)
}
