package domain

import (
	"context"
)

type ShopProductCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Category    string `json:"category"`
	ImageURL    string `json:"image_url"`
}

type ProductStock struct {
	WarehouseID int64 `json:"warehouse_id"`
	Stock       int64 `json:"stock"`
}

type ProductCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Category    string `json:"category"`
	ImageURL    string `json:"image_url"`
	ShopID      int64  `json:"shop_id"`
}

type ProductCreateResponse struct {
	ID int64 `json:"id"`
}

type ProductUsecase interface {
	CreateProduct(ctx context.Context, shopID int64, req ShopProductCreateRequest) (ProductCreateResponse, error)
}

type ProductRepository interface {
	CreateProduct(ctx context.Context, req ProductCreateRequest) (ProductCreateResponse, error)
}
