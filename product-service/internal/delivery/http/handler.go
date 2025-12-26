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

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.ProdUsecase.GetAllProducts(r.Context())
	if err != nil {
		h.respondWithError(w, pkgerrors.GetStatusCode(err), err.Error())
		return
	}
	h.respondWithJSON(w, http.StatusOK, products)
}

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
