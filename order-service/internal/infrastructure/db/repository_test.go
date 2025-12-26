package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/user/go-microservices/order-service/internal/domain"
	"github.com/user/go-microservices/pkg/logger"
	"github.com/user/go-microservices/pkg/valueobject"
)

func TestOrderRepository(t *testing.T) {
	logger.Init()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	repo := NewOrderRepository(db)

	t.Run("Create_Success", func(t *testing.T) {
		order := &domain.Order{
			UserID:     1,
			ProductID:  1,
			UnitPrice:  valueobject.NewMoney(100),
			Quantity:   1,
			TotalPrice: valueobject.NewMoney(100),
		}

		mock.ExpectQuery("INSERT INTO orders").
			WithArgs(order.UserID, order.ProductID, sqlmock.AnyArg(), order.UnitPrice.Amount(), order.Quantity, order.TotalPrice.Amount(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		err := repo.Create(context.Background(), order)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), order.ID)
	})

	t.Run("GetByID_Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "product_id", "product_name", "unit_price", "quantity", "total_price", "order_status", "payment_status", "created_at"}).
			AddRow(1, 1, 1, "Product 1", 100.0, 1, 100.0, "PENDING", "PENDING", time.Now())

		mock.ExpectQuery("SELECT (.+) FROM orders WHERE id = \\$1").
			WithArgs(int64(1)).
			WillReturnRows(rows)

		order, err := repo.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, int64(1), order.ID)
	})
}
