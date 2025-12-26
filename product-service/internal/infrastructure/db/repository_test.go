package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/user/go-microservices/pkg/logger"
	"github.com/user/go-microservices/pkg/valueobject"
	"github.com/user/go-microservices/product-service/internal/domain"
)

func TestProductRepository(t *testing.T) {
	logger.Init()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	repo := NewPostgresRepository(db)

	t.Run("Create_Success", func(t *testing.T) {
		p := &domain.Product{
			SKU:   "SKU1",
			Name:  "Product 1",
			Price: valueobject.NewMoney(100),
		}

		mock.ExpectQuery("INSERT INTO products").
			WithArgs(p.SKU, p.Name, sqlmock.AnyArg(), p.Price.Amount(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		err := repo.Create(context.Background(), p)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), p.ID)
	})

	t.Run("ReserveStock_Success", func(t *testing.T) {
		mock.ExpectExec("UPDATE products SET reserved_qty = reserved_qty \\+ \\$1").
			WithArgs(5, int64(1)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.ReserveStock(context.Background(), 1, 5)

		assert.NoError(t, err)
	})
}
