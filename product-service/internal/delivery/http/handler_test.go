package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/user/go-microservices/pkg/logger"
	"github.com/user/go-microservices/product-service/internal/domain"
	"github.com/user/go-microservices/product-service/internal/usecase/mocks"
)

func TestProductHandler(t *testing.T) {
	logger.Init()
	mockUC := mocks.NewProductUsecase(t)
	router := mux.NewRouter()
	NewProductHandler(router, mockUC)

	t.Run("CreateProduct_Success", func(t *testing.T) {
		p := domain.Product{SKU: "SKU1", Name: "Product 1"}
		body, _ := json.Marshal(p)
		req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		mockUC.On("CreateProduct", mock.Anything, mock.AnythingOfType("*domain.Product")).Return(nil)

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	t.Run("GetProduct_Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/products/1", nil)
		rr := httptest.NewRecorder()

		mockUC.On("GetProduct", mock.Anything, int64(1)).Return(&domain.Product{ID: 1, SKU: "SKU1"}, nil)

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var res domain.Product
		json.Unmarshal(rr.Body.Bytes(), &res)
		assert.Equal(t, int64(1), res.ID)
	})
}
