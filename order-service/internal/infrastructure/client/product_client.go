package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/go-microservices/order-service/internal/domain"
	pkgerrors "github.com/user/go-microservices/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type productClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewProductClient(baseURL string) domain.ProductClient {
	return &productClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
			Timeout:   5 * time.Second,
		},
	}
}

func (c *productClient) GetProduct(ctx context.Context, id int64) (*domain.ProductView, error) {
	url := fmt.Sprintf("%s/products/%d", c.baseURL, id)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, pkgerrors.ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, pkgerrors.ErrInternal
	}

	var p domain.ProductView
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

type stockReq struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

func (c *productClient) ReserveStock(ctx context.Context, id int64, qty int) error {
	url := fmt.Sprintf("%s/products/reserve", c.baseURL)
	body, _ := json.Marshal(stockReq{ProductID: id, Quantity: qty})

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnprocessableEntity {
		return pkgerrors.ErrInsufficientStock
	}
	if resp.StatusCode != http.StatusOK {
		return pkgerrors.ErrInternal
	}
	return nil
}

func (c *productClient) ReleaseStock(ctx context.Context, id int64, qty int) error {
	url := fmt.Sprintf("%s/products/release", c.baseURL)
	body, _ := json.Marshal(stockReq{ProductID: id, Quantity: qty})

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return pkgerrors.ErrInternal
	}
	return nil
}
