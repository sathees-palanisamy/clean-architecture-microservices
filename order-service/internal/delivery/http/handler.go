package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/user/go-microservices/order-service/internal/domain"
	"github.com/user/go-microservices/order-service/internal/usecase"
	pkgerrors "github.com/user/go-microservices/pkg/errors"
	"github.com/user/go-microservices/pkg/logger"
	"go.uber.org/zap"
)

var _ = domain.Order{}

type OrderHandler struct {
	OrderUsecase usecase.OrderUsecase
}

func NewOrderHandler(r *mux.Router, us usecase.OrderUsecase) {
	handler := &OrderHandler{
		OrderUsecase: us,
	}

	r.HandleFunc("/orders", handler.CreateOrder).Methods("POST")
	r.HandleFunc("/orders", handler.GetAllOrders).Methods("GET")
	r.HandleFunc("/orders/{id}", handler.GetOrder).Methods("GET")
	r.HandleFunc("/health", handler.HealthCheck).Methods("GET")
}

type CreateOrderRequest struct {
	UserID    int64 `json:"user_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order for a product and user
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order body CreateOrderRequest true "Order request"
// @Success 201 {object} domain.Order
// @Failure 400 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Quantity <= 0 {
		h.respondWithError(w, http.StatusBadRequest, "Quantity must be greater than 0")
		return
	}

	ctx := r.Context()
	order, err := h.OrderUsecase.CreateOrder(ctx, req.UserID, req.ProductID, req.Quantity)
	if err != nil {
		h.respondWithError(w, pkgerrors.GetStatusCode(err), err.Error())
		return
	}
	h.respondWithJSON(w, http.StatusCreated, order)
}

// GetAllOrders godoc
// @Summary List all orders
// @Description Get a list of all orders
// @Tags orders
// @Produce  json
// @Success 200 {array} domain.Order
// @Router /orders [get]
func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.OrderUsecase.GetAllOrders(r.Context())
	if err != nil {
		h.respondWithError(w, pkgerrors.GetStatusCode(err), err.Error())
		return
	}
	h.respondWithJSON(w, http.StatusOK, orders)
}

// GetOrder godoc
// @Summary Get an order by ID
// @Description Get detailed information about an order by its ID
// @Tags orders
// @Produce  json
// @Param id path int true "Order ID"
// @Success 200 {object} domain.Order
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	o, err := h.OrderUsecase.GetOrder(r.Context(), id)
	if err != nil {
		h.respondWithError(w, pkgerrors.GetStatusCode(err), err.Error())
		return
	}
	h.respondWithJSON(w, http.StatusOK, o)
}

func (h *OrderHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.respondWithJSON(w, http.StatusOK, map[string]string{"status": "UP"})
}

func (h *OrderHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

func (h *OrderHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)

	logger.Info("request handled",
		zap.Int("status", code),
		zap.String("response", string(response)),
	)
}
