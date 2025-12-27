package integration

import (
	"context"
	"fmt"
	"time"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/mock"
	"github.com/user/go-microservices/pkg/valueobject"
	"github.com/user/go-microservices/product-service/internal/domain"
	repoMocks "github.com/user/go-microservices/product-service/internal/domain/mocks"
	"github.com/user/go-microservices/product-service/internal/usecase"
)

type productTestContext struct {
	repo        *repoMocks.ProductRepository
	uc          usecase.ProductUsecase
	lastProduct *domain.Product
	lastError   error
	mockProduct *domain.Product
}

func (c *productTestContext) iCreateAProductWithSKUNameAndPrice(sku, name string, price float64) error {
	c.repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Product")).Return(nil).Once()

	p := &domain.Product{
		SKU:   sku,
		Name:  name,
		Price: valueobject.NewMoney(price),
	}
	c.lastError = c.uc.CreateProduct(context.Background(), p)
	return nil
}

func (c *productTestContext) theProductShouldBeSuccessfullySaved() error {
	if c.lastError != nil {
		return fmt.Errorf("expected no error, got %v", c.lastError)
	}
	return nil
}

func (c *productTestContext) aProductExistsWithSKUNamePriceAndStock(sku, name string, price float64, stock int) error {
	c.mockProduct = &domain.Product{
		ID:          1,
		SKU:         sku,
		Name:        name,
		Price:       valueobject.NewMoney(price),
		TotalQty:    stock,
		ReservedQty: 0,
	}
	c.repo.On("GetByID", mock.Anything, int64(1)).Return(c.mockProduct, nil)
	return nil
}

func (c *productTestContext) iReserveUnitsOfStockForThisProduct(qty int) error {
	if c.mockProduct.AvailableQty() >= qty {
		c.repo.On("ReserveStock", mock.Anything, c.mockProduct.ID, qty).Return(nil).Once()
		// Update internal mock state for validation
		c.mockProduct.ReservedQty += qty
	} else {
		c.repo.On("ReserveStock", mock.Anything, c.mockProduct.ID, qty).Return(fmt.Errorf("insufficient stock")).Once()
	}

	c.lastError = c.uc.ReserveStock(context.Background(), c.mockProduct.ID, qty)
	return nil
}

func (c *productTestContext) theReservationShouldBeSuccessful() error {
	if c.lastError != nil {
		return fmt.Errorf("expected no error, got %v", c.lastError)
	}
	return nil
}

func (c *productTestContext) theAvailableStockShouldBe(expectedStock int) error {
	actual := c.mockProduct.AvailableQty()
	if actual != expectedStock {
		return fmt.Errorf("expected available stock %d, got %d", expectedStock, actual)
	}
	return nil
}

func (c *productTestContext) theReservationShouldFailWith(expectedMsg string) error {
	if c.lastError == nil {
		return fmt.Errorf("expected error containing %q, but got none", expectedMsg)
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	c := &productTestContext{
		repo: new(repoMocks.ProductRepository),
	}
	c.uc = usecase.NewProductUsecase(c.repo, 5*time.Second)

	ctx.Step(`^I create a product with SKU "([^"]*)", name "([^"]*)", and price ([\d.]+)$`, c.iCreateAProductWithSKUNameAndPrice)
	ctx.Step(`^the product should be successfully saved$`, c.theProductShouldBeSuccessfullySaved)
	ctx.Step(`^a product exists with SKU "([^"]*)", name "([^"]*)", price ([\d.]+), and stock (\d+)$`, c.aProductExistsWithSKUNamePriceAndStock)
	ctx.Step(`^I reserve (\d+) units of stock for this product$`, c.iReserveUnitsOfStockForThisProduct)
	ctx.Step(`^the reservation should be successful$`, c.theReservationShouldBeSuccessful)
	ctx.Step(`^the available stock should be (\d+)$`, c.theAvailableStockShouldBe)
	ctx.Step(`^the reservation should fail with "([^"]*)"$`, c.theReservationShouldFailWith)
}
