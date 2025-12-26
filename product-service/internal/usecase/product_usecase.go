package usecase

import (
	"context"
	"time"

	"github.com/user/go-microservices/product-service/internal/domain"
)

type ProductUsecase interface {
	CreateProduct(ctx context.Context, p *domain.Product) error
	GetProduct(ctx context.Context, id int64) (*domain.Product, error)
	ReserveStock(ctx context.Context, id int64, qty int) error
	ReleaseStock(ctx context.Context, id int64, qty int) error
	ConfirmStock(ctx context.Context, id int64, qty int) error
	GetAllProducts(ctx context.Context) ([]*domain.Product, error)
}

type productUsecase struct {
	repo           domain.ProductRepository
	contextTimeout time.Duration
}

func NewProductUsecase(repo domain.ProductRepository, timeout time.Duration) ProductUsecase {
	return &productUsecase{
		repo:           repo,
		contextTimeout: timeout,
	}
}

func (u *productUsecase) CreateProduct(ctx context.Context, p *domain.Product) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.Create(ctx, p)
}

func (u *productUsecase) GetProduct(ctx context.Context, id int64) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.GetByID(ctx, id)
}

func (u *productUsecase) ReserveStock(ctx context.Context, id int64, qty int) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.ReserveStock(ctx, id, qty)
}

func (u *productUsecase) ReleaseStock(ctx context.Context, id int64, qty int) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.ReleaseStock(ctx, id, qty)
}

func (u *productUsecase) ConfirmStock(ctx context.Context, id int64, qty int) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.ConfirmStock(ctx, id, qty)
}

func (u *productUsecase) GetAllProducts(ctx context.Context) ([]*domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.GetAll(ctx)
}
