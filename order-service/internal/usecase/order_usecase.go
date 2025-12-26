package usecase

import (
	"context"
	"time"

	"github.com/user/go-microservices/order-service/internal/domain"
	pkgerrors "github.com/user/go-microservices/pkg/errors"
	"github.com/user/go-microservices/pkg/logger"
	"go.uber.org/zap"
)

type OrderUsecase interface {
	CreateOrder(ctx context.Context, userID, productID int64, qty int) (*domain.Order, error)
	GetOrder(ctx context.Context, id int64) (*domain.Order, error)
	GetAllOrders(ctx context.Context) ([]*domain.Order, error)
}

type orderUsecase struct {
	repo           domain.OrderRepository
	productClient  domain.ProductClient
	contextTimeout time.Duration
}

func NewOrderUsecase(repo domain.OrderRepository, pClient domain.ProductClient, timeout time.Duration) OrderUsecase {
	return &orderUsecase{
		repo:           repo,
		productClient:  pClient,
		contextTimeout: timeout,
	}
}

func (u *orderUsecase) CreateOrder(ctx context.Context, userID, productID int64, qty int) (*domain.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// 1. Get Product Details (Snapshot)
	product, err := u.productClient.GetProduct(ctx, productID)
	if err != nil {
		return nil, err
	}

	// 2. Reserve Stock
	err = u.productClient.ReserveStock(ctx, productID, qty)
	if err != nil {
		return nil, err
	}

	// 3. Create Order Aggregate
	order, err := domain.NewOrder(userID, productID, product.Name, product.Price, qty)
	if err != nil {
		// Rollback: Release Stock
		logger.FromContext(ctx).Warn("invalid order parameters, rolling back stock", zap.Int64("product_id", productID))
		if rbErr := u.productClient.ReleaseStock(context.Background(), productID, qty); rbErr != nil {
			logger.FromContext(ctx).Error("failed to rollback stock", zap.Error(rbErr))
		}
		return nil, err
	}

	if err := u.repo.Create(ctx, order); err != nil {
		// Rollback: Release Stock
		logger.FromContext(ctx).Warn("rolling back stock reservation due to order creation failure", zap.Int64("product_id", productID))
		if rbErr := u.productClient.ReleaseStock(context.Background(), productID, qty); rbErr != nil {
			logger.FromContext(ctx).Error("failed to rollback stock", zap.Error(rbErr))
		}
		return nil, pkgerrors.ErrInternal
	}

	return order, nil
}

func (u *orderUsecase) GetOrder(ctx context.Context, id int64) (*domain.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.GetByID(ctx, id)
}

func (u *orderUsecase) GetAllOrders(ctx context.Context) ([]*domain.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.GetAll(ctx)
}
