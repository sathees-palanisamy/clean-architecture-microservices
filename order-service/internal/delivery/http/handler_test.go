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
	"github.com/user/go-microservices/order-service/internal/domain"
	"github.com/user/go-microservices/order-service/internal/usecase/mocks"
	"github.com/user/go-microservices/pkg/logger"
)

func TestOrderHandler(t *testing.T) {
	logger.Init()
	mockUC := mocks.NewOrderUsecase(t)
	router := mux.NewRouter()
	NewOrderHandler(router, mockUC)

	t.Run("CreateOrder_Success", func(t *testing.T) {
		reqBody := CreateOrderRequest{UserID: 1, ProductID: 1, Quantity: 2}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		mockUC.On("CreateOrder", mock.Anything, int64(1), int64(1), 2).Return(&domain.Order{ID: 1}, nil)

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		var res domain.Order
		json.Unmarshal(rr.Body.Bytes(), &res)
		assert.Equal(t, int64(1), res.ID)
	})

	t.Run("GetOrder_Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/orders/1", nil)
		rr := httptest.NewRecorder()

		mockUC.On("GetOrder", mock.Anything, int64(1)).Return(&domain.Order{ID: 1}, nil)

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var res domain.Order
		json.Unmarshal(rr.Body.Bytes(), &res)
		assert.Equal(t, int64(1), res.ID)
	})
}
