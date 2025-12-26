package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/user/go-microservices/order-service/internal/domain"
	"github.com/user/go-microservices/order-service/internal/domain/mocks"
	"github.com/user/go-microservices/pkg/logger"
	"github.com/user/go-microservices/pkg/valueobject"
)

func TestOrderUsecase_CreateOrder(t *testing.T) {
	logger.Init()
	timeout := 5 * time.Second

	t.Run("Success", func(t *testing.T) {
		mockRepo := mocks.NewOrderRepository(t)
		mockProductClient := mocks.NewProductClient(t)
		uc := NewOrderUsecase(mockRepo, mockProductClient, timeout)

		product := &domain.ProductView{
			ID:    1,
			Name:  "Test Product",
			Price: valueobject.NewMoney(100.0),
		}

		mockProductClient.On("GetProduct", mock.Anything, int64(1)).Return(product, nil)
		mockProductClient.On("ReserveStock", mock.Anything, int64(1), 2).Return(nil)
		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Order")).Return(nil)

		order, err := uc.CreateOrder(context.Background(), 101, 1, 2)

		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, int64(101), order.UserID)
		assert.Equal(t, valueobject.NewMoney(200.0), order.TotalPrice)
	})

	t.Run("ProductNotFound", func(t *testing.T) {
		mockRepo := mocks.NewOrderRepository(t)
		mockProductClient := mocks.NewProductClient(t)
		uc := NewOrderUsecase(mockRepo, mockProductClient, timeout)

		mockProductClient.On("GetProduct", mock.Anything, int64(2)).Return(nil, assert.AnError)

		order, err := uc.CreateOrder(context.Background(), 101, 2, 1)

		assert.Error(t, err)
		assert.Nil(t, order)
	})

	t.Run("InsufficientStock", func(t *testing.T) {
		mockRepo := mocks.NewOrderRepository(t)
		mockProductClient := mocks.NewProductClient(t)
		uc := NewOrderUsecase(mockRepo, mockProductClient, timeout)

		product := &domain.ProductView{
			ID:    1,
			Name:  "Test Product",
			Price: valueobject.NewMoney(100.0),
		}
		mockProductClient.On("GetProduct", mock.Anything, int64(1)).Return(product, nil)
		mockProductClient.On("ReserveStock", mock.Anything, int64(1), 10).Return(assert.AnError)

		order, err := uc.CreateOrder(context.Background(), 101, 1, 10)

		assert.Error(t, err)
		assert.Nil(t, order)
	})

	t.Run("RepoFailure_WithRollback", func(t *testing.T) {
		mockRepo := mocks.NewOrderRepository(t)
		mockProductClient := mocks.NewProductClient(t)
		uc := NewOrderUsecase(mockRepo, mockProductClient, timeout)

		product := &domain.ProductView{
			ID:    1,
			Name:  "Test Product",
			Price: valueobject.NewMoney(100.0),
		}
		mockProductClient.On("GetProduct", mock.Anything, int64(1)).Return(product, nil)
		mockProductClient.On("ReserveStock", mock.Anything, int64(1), 1).Return(nil)
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(assert.AnError)
		mockProductClient.On("ReleaseStock", mock.Anything, int64(1), 1).Return(nil)

		order, err := uc.CreateOrder(context.Background(), 101, 1, 1)

		assert.Error(t, err)
		assert.Nil(t, order)
	})
}
