package domain

import "context"

type InitStockRequest struct {
	ShopID    int64 `json:"shop_id"`
	ProductID int64 `json:"product_id"`
}

type StockRepository interface {
	InitStock(ctx context.Context, req InitStockRequest) error
}

type StockUsecase interface {
	InitStock(ctx context.Context, req InitStockRequest) error
}
