package stockrepo

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

type stockRepository struct {
	httpClient         *http.Client
	baseURL            string
	internalAuthHeader string
}

func NewStockRepository(httpClient *http.Client, baseURL string, internalAuthHeader string) domain.StockRepository {
	return &stockRepository{
		httpClient:         httpClient,
		baseURL:            baseURL,
		internalAuthHeader: internalAuthHeader,
	}
}

func (r *stockRepository) InitStock(ctx context.Context, warehouse domain.InitStockRequest) error {
	url := fmt.Sprintf("%s/internal/warehouse-service/stocks", r.baseURL)
	reqBody, err := json.Marshal(warehouse)
	if err != nil {
		slog.ErrorContext(ctx, "[stockRepository] InitStock", "json.Marshal", err)
		return err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		slog.ErrorContext(ctx, "[stockRepository] InitStock", "http.NewRequestWithContext", err)
		return err
	}

	pkg.AddRequestHeader(ctx, r.internalAuthHeader, httpReq)

	resp, err := r.httpClient.Do(httpReq)
	if err != nil {
		slog.ErrorContext(ctx, "[stockRepository] InitStock", "httpClient.Do", err)
		return err
	}
	defer resp.Body.Close()

	var res any
	if err := pkg.DecodeResponseBody(resp, &res); err != nil {
		slog.ErrorContext(ctx, "[stockRepository] InitStock", "DecodeResponseBody", err)
		return err
	}

	return nil
}
