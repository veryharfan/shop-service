package domain

import (
	"context"
)

type UserShopUpdateRequest struct {
	ShopID int64 `json:"shop_id"`
}

type UserRepository interface {
	PatchUserShop(ctx context.Context, userID int64, req UserShopUpdateRequest) error
}
