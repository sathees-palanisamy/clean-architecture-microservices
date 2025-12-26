package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/user/go-microservices/pkg/logger"
	"github.com/user/go-microservices/pkg/valueobject"
	"github.com/user/go-microservices/product-service/internal/domain"
	"github.com/user/go-microservices/product-service/internal/domain/mocks"
)

func TestProductUsecase(t *testing.T) {
	logger.Init()
	mockRepo := mocks.NewProductRepository(t)
	timeout := 5 * time.Second
	uc := NewProductUsecase(mockRepo, timeout)

	ctx := context.Background()

	t.Run("CreateProduct", func(t *testing.T) {
		p := &domain.Product{SKU: "SKU1", Name: "N1", Price: valueobject.NewMoney(10)}
		mockRepo.On("Create", mock.Anything, p).Return(nil).Once()
		err := uc.CreateProduct(ctx, p)
		assert.NoError(t, err)
	})

	t.Run("GetProduct", func(t *testing.T) {
		p := &domain.Product{ID: 1, SKU: "SKU1"}
		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(p, nil).Once()
		res, err := uc.GetProduct(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, p, res)
	})

	t.Run("ReserveStock", func(t *testing.T) {
		mockRepo.On("ReserveStock", mock.Anything, int64(1), 5).Return(nil).Once()
		err := uc.ReserveStock(ctx, 1, 5)
		assert.NoError(t, err)
	})
}
