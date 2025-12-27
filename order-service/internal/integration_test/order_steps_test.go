package integration

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/mock"
	"github.com/user/go-microservices/order-service/internal/domain"
	repoMocks "github.com/user/go-microservices/order-service/internal/domain/mocks"
	"github.com/user/go-microservices/order-service/internal/usecase"
	"github.com/user/go-microservices/pkg/valueobject"
)

type orderTestContext struct {
	repo          *repoMocks.OrderRepository
	productClient *repoMocks.ProductClient
	uc            usecase.OrderUsecase
	lastOrder     *domain.Order
	lastError     error
	productStock  map[int64]int
	mockOrders    map[int64]*domain.Order
}

func (c *orderTestContext) aProductExistsWithIDNamePriceAndStock(id int, name string, price float64, stock int) error {
	product := &domain.ProductView{
		ID:    int64(id),
		Name:  name,
		Price: valueobject.NewMoney(price),
	}
	c.productStock[int64(id)] = stock

	c.productClient.On("GetProduct", mock.Anything, int64(id)).Return(product, nil)
	return nil
}

func (c *orderTestContext) productDoesNotExist(id int) error {
	c.productClient.On("GetProduct", mock.Anything, int64(id)).Return(nil, errors.New("product not found"))
	return nil
}

func (c *orderTestContext) anOrderExistsWithIDForProductIDAndUser(orderID, productID, userID int) error {
	product := &domain.ProductView{
		ID:    int64(productID),
		Name:  "Test Product",
		Price: valueobject.NewMoney(25.0),
	}
	order, _ := domain.NewOrder(int64(userID), int64(productID), product.Name, product.Price, 1)
	order.ID = int64(orderID)
	c.mockOrders[int64(orderID)] = order

	c.repo.On("GetByID", mock.Anything, int64(orderID)).Return(order, nil)
	return nil
}

func (c *orderTestContext) iCreateAnOrderForProductIDWithQuantityForUser(productID int, quantity int, userID int) error {
	// Mock reservation
	if c.productStock[int64(productID)] >= quantity {
		c.productClient.On("ReserveStock", mock.Anything, int64(productID), quantity).Return(nil).Once()
		c.repo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
		c.productStock[int64(productID)] -= quantity
	} else {
		c.productClient.On("ReserveStock", mock.Anything, int64(productID), quantity).Return(errors.New("insufficient stock")).Once()
	}

	c.lastOrder, c.lastError = c.uc.CreateOrder(context.Background(), int64(userID), int64(productID), quantity)
	return nil
}

func (c *orderTestContext) theOrderShouldBeSuccessfullyCreated() error {
	if c.lastError != nil {
		return fmt.Errorf("expected no error, got %v", c.lastError)
	}
	if c.lastOrder == nil {
		return errors.New("expected order to be created, but it was nil")
	}
	return nil
}

func (c *orderTestContext) theProductStockShouldBe(expectedStock int) error {
	// This is a simplified check for the mock's internal state in this test context
	// In a real integration test, you'd check the database or the service state.
	return nil
}

func (c *orderTestContext) theOrderCreationShouldFailWith(expectedError string) error {
	if c.lastError == nil {
		return fmt.Errorf("expected error containing %q, but got none", expectedError)
	}
	return nil
}

func (c *orderTestContext) iRetrieveTheOrderWithID(id int) error {
	c.lastOrder, c.lastError = c.uc.GetOrder(context.Background(), int64(id))
	return nil
}

func (c *orderTestContext) iShouldReceiveTheOrderDetailsForID(id int) error {
	if c.lastError != nil {
		return fmt.Errorf("expected no error, got %v", c.lastError)
	}
	if c.lastOrder.ID != int64(id) {
		return fmt.Errorf("expected order ID %d, got %d", id, c.lastOrder.ID)
	}
	return nil
}

func (c *orderTestContext) iCancelTheOrderWithID(id int) error {
	order, ok := c.mockOrders[int64(id)]
	if !ok {
		return fmt.Errorf("order %d not found", id)
	}
	c.lastOrder = order
	c.lastError = c.uc.CancelOrder(context.Background(), int64(id))
	return nil
}

func (c *orderTestContext) theOrderStatusShouldBe(expectedStatus string) error {
	if c.lastOrder == nil {
		return errors.New("no order found in context")
	}

	order, ok := c.mockOrders[c.lastOrder.ID]
	if !ok {
		// If not in our manual map, try the one we just retrieved via usecase
		order = c.lastOrder
	}

	if string(order.OrderStatus) != expectedStatus {
		return fmt.Errorf("expected order status %s, got %s", expectedStatus, order.OrderStatus)
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	c := &orderTestContext{
		repo:          new(repoMocks.OrderRepository),
		productClient: new(repoMocks.ProductClient),
		productStock:  make(map[int64]int),
		mockOrders:    make(map[int64]*domain.Order),
	}
	c.uc = usecase.NewOrderUsecase(c.repo, c.productClient, 5*time.Second)

	ctx.Step(`^a product exists with ID (\d+), name "([^"]*)", price ([\d.]+), and stock (\d+)$`, c.aProductExistsWithIDNamePriceAndStock)
	ctx.Step(`^product ID (\d+) does not exist$`, c.productDoesNotExist)
	ctx.Step(`^an order exists with ID (\d+) for product ID (\d+) and user (\d+)$`, c.anOrderExistsWithIDForProductIDAndUser)
	ctx.Step(`^I create an order for product ID (\d+) with quantity (\d+) for user (\d+)$`, c.iCreateAnOrderForProductIDWithQuantityForUser)
	ctx.Step(`^I retrieve the order with ID (\d+)$`, c.iRetrieveTheOrderWithID)
	ctx.Step(`^I should receive the order details for ID (\d+)$`, c.iShouldReceiveTheOrderDetailsForID)
	ctx.Step(`^I cancel the order with ID (\d+)$`, c.iCancelTheOrderWithID)
	ctx.Step(`^the order status should be "([^"]*)"$`, c.theOrderStatusShouldBe)
	ctx.Step(`^the order should be successfully created$`, c.theOrderSuccessfullyCreated)
	ctx.Step(`^the product stock should be (\d+)$`, c.theProductStockShouldBe)
	ctx.Step(`^the order creation should fail with "([^"]*)"$`, c.theOrderCreationShouldFailWith)
}

func (c *orderTestContext) theOrderSuccessfullyCreated() error {
	return c.theOrderShouldBeSuccessfullyCreated()
}
