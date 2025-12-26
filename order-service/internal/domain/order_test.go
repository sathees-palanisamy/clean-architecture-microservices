package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/user/go-microservices/pkg/valueobject"
)

func TestOrder_Aggregate(t *testing.T) {
	price := valueobject.NewMoney(100)

	t.Run("NewOrder_Valid", func(t *testing.T) {
		order, err := NewOrder(1, 1, "Test", price, 2)
		assert.NoError(t, err)
		assert.Equal(t, valueobject.NewMoney(200), order.TotalPrice)
		assert.Equal(t, OrderPending, order.OrderStatus)
	})

	t.Run("NewOrder_InvalidQty", func(t *testing.T) {
		order, err := NewOrder(1, 1, "Test", price, 0)
		assert.Error(t, err)
		assert.Nil(t, order)
	})

	t.Run("StateTransitions", func(t *testing.T) {
		order, _ := NewOrder(1, 1, "Test", price, 1)

		// Pay
		err := order.Pay()
		assert.NoError(t, err)
		assert.Equal(t, PaymentPaid, order.PaymentStatus)

		// Complete
		err = order.Complete()
		assert.NoError(t, err)
		assert.Equal(t, OrderCompleted, order.OrderStatus)

		// Cancel after complete (Fail)
		err = order.Cancel()
		assert.Error(t, err)
	})

	t.Run("Cancel_Valid", func(t *testing.T) {
		order, _ := NewOrder(1, 1, "Test", price, 1)
		err := order.Cancel()
		assert.NoError(t, err)
		assert.Equal(t, OrderCancelled, order.OrderStatus)

		// Pay after cancel (Fail)
		err = order.Pay()
		assert.Error(t, err)
	})
}
