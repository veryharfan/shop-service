package productrepo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"shop-service/app/domain"
	"shop-service/pkg"
)

type productRepository struct {
	httpClient         *http.Client
	baseURL            string
	internalAuthHeader string
}

func NewProductRepository(httpClient *http.Client, baseURL string, internalAuthHeader string) domain.ProductRepository {
	return &productRepository{
		httpClient:         httpClient,
		baseURL:            baseURL,
		internalAuthHeader: internalAuthHeader,
	}
}

func (r *productRepository) CreateProduct(ctx context.Context, req domain.ProductCreateRequest) (domain.ProductCreateResponse, error) {
	url := fmt.Sprintf("%s/internal/product-service/products", r.baseURL)
	reqBody, err := json.Marshal(req)
	if err != nil {
		slog.ErrorContext(ctx, "[productRepository] CreateProduct", "json.Marshal", err)
		return domain.ProductCreateResponse{}, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		slog.ErrorContext(ctx, "[productRepository] CreateProduct", "http.NewRequestWithContext", err)
		return domain.ProductCreateResponse{}, err
	}

	pkg.AddRequestHeader(ctx, r.internalAuthHeader, httpReq)

	resp, err := r.httpClient.Do(httpReq)
	if err != nil {
		slog.ErrorContext(ctx, "[productRepository] CreateProduct", "httpClient.Do", err)
		return domain.ProductCreateResponse{}, err
	}
	defer resp.Body.Close()

	var response domain.ProductCreateResponse
	if err := pkg.DecodeResponseBody(resp, &response); err != nil {
		slog.ErrorContext(ctx, "[productRepository] CreateProduct", "DecodeResponseBody", err)
		return domain.ProductCreateResponse{}, err
	}

	return response, nil
}
