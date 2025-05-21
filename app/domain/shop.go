package domain

import (
	"context"
	"database/sql"
)

type Shop struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateShopRequest struct {
	UserID int64  `json:"user_id" validate:"required"`
	Name   string `json:"name" validate:"required"`
}

type ShopUsecase interface {
	Create(ctx context.Context, req CreateShopRequest) (*Shop, error)
	GetByUserID(ctx context.Context, userID int64) (*Shop, error)
}

type ShopRepository interface {
	Create(ctx context.Context, req *Shop, tx *sql.Tx) error
	GetByUserID(ctx context.Context, userID int64) (*Shop, error)

	WithTransaction(ctx context.Context, fn func(context.Context, *sql.Tx) error) error
}
