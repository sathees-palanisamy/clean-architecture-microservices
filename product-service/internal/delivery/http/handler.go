package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	pkgerrors "github.com/user/go-microservices/pkg/errors"
	"github.com/user/go-microservices/pkg/logger"
	"github.com/user/go-microservices/product-service/internal/domain"
	"github.com/user/go-microservices/product-service/internal/usecase"
	"go.uber.org/zap"
)

type ProductHandler struct {
	ProdUsecase usecase.ProductUsecase
}

func NewProductHandler(r *mux.Router, us usecase.ProductUsecase) {
	handler := &ProductHandler{
		ProdUsecase: us,
	}

	r.HandleFunc("/products", handler.CreateProduct).Methods("POST")
	r.HandleFunc("/products", handler.GetAllProducts).Methods("GET")
	r.HandleFunc("/products/{id}", handler.GetProduct).Methods("GET")
	r.HandleFunc("/products/reserve", handler.ReserveStock).Methods("POST")
	r.HandleFunc("/products/release", handler.ReleaseStock).Methods("POST")
	r.HandleFunc("/products/confirm", handler.ConfirmStock).Methods("POST")
	r.HandleFunc("/health", handler.HealthCheck).Methods("GET")
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with the provided details
// @Tags products
// @Accept  json
// @Produce  json
// @Param product body domain.Product true "Product object"
// @Success 201 {object} domain.Product
// @Failure 400 {object} map[string]string
// @Router /products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p domain.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	// Basic validation
	if p.TotalQty < 0 || p.Price.IsNegative() {
		h.respondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	ctx := r.Context()
	err := h.ProdUsecase.CreateProduct(ctx, &p)
	if err != nil {
		h.respondWithError(w, pkgerrors.GetStatusCode(err), err.Error())
		return
	}
	h.respondWithJSON(w, http.StatusCreated, p)
}

// GetAllProducts godoc
// @Summary List all products
// @Description Get a list of all products
// @Tags products
// @Produce  json
// @Success 200 {array} domain.Product
// @Router /products [get]
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.ProdUsecase.GetAllProducts(r.Context())
	if err != nil {
		h.respondWithError(w, pkgerrors.GetStatusCode(err), err.Error())
		return
	}
	h.respondWithJSON(w, http.StatusOK, products)
}

// GetProduct godoc
// @Summary Get a product by ID
// @Description Get detailed information about a product by its ID
// @Tags products
// @Produce  json
// @Param id path int true "Product ID"
// @Success 200 {object} domain.Product
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	p, err := h.ProdUsecase.GetProduct(r.Context(), id)
	if err != nil {
		h.respondWithError(w, pkgerrors.GetStatusCode(err), err.Error())
		return
	}
	h.respondWithJSON(w, http.StatusOK, p)
}

type StockRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

// ReserveStock godoc
// @Summary Reserve stock for a product
// @Description Reserve a specific quantity of stock for a given product ID
// @Tags stock
// @Accept  json
// @Produce  json
// @Param request body StockRequest true "Stock reservation request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Router /products/reserve [post]
func (h *ProductHandler) ReserveStock(w http.ResponseWriter, r *http.Request) {
	var req StockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	err := h.ProdUsecase.ReserveStock(r.Context(), req.ProductID, req.Quantity)
	if err != nil {
		h.respondWithError(w, pkgerrors.GetStatusCode(err), err.Error())
		return
	}
	h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "reserved"})
}

// ReleaseStock godoc
// @Summary Release reserved stock
// @Description Release a previously reserved quantity of stock
// @Tags stock
// @Accept  json
// @Produce  json
// @Param request body StockRequest true "Stock release request"
// @Success 200 {object} map[string]string
// @Router /products/release [post]
func (h *ProductHandler) ReleaseStock(w http.ResponseWriter, r *http.Request) {
	var req StockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	err := h.ProdUsecase.ReleaseStock(r.Context(), req.ProductID, req.Quantity)
	if err != nil {
		h.respondWithError(w, pkgerrors.GetStatusCode(err), err.Error())
		return
	}
	h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "released"})
}

// ConfirmStock godoc
// @Summary Confirm stock reservation
// @Description Confirm the reservation and permanently deduct stock
// @Tags stock
// @Accept  json
// @Produce  json
// @Param request body StockRequest true "Stock confirmation request"
// @Success 200 {object} map[string]string
// @Router /products/confirm [post]
func (h *ProductHandler) ConfirmStock(w http.ResponseWriter, r *http.Request) {
	var req StockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	err := h.ProdUsecase.ConfirmStock(r.Context(), req.ProductID, req.Quantity)
	if err != nil {
		h.respondWithError(w, pkgerrors.GetStatusCode(err), err.Error())
		return
	}
	h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "confirmed"})
}

func (h *ProductHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.respondWithJSON(w, http.StatusOK, map[string]string{"status": "UP"})
}

func (h *ProductHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

func (h *ProductHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)

	// Structured logging for request
	logger.Info("request handled",
		zap.Int("status", code),
		zap.String("response", string(response)),
	)
}
