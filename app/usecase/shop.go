package usecase

import (
	"context"
	"log/slog"
	"shop-service/app/domain"
	"shop-service/config"
)

type shopUsecase struct {
	shopRepo domain.ShopRepository
	cfg      *config.Config
}

func NewShopUsecase(shopRepo domain.ShopRepository, cfg *config.Config) domain.ShopUsecase {
	return &shopUsecase{shopRepo, cfg}
}

func (u *shopUsecase) Create(ctx context.Context, req domain.CreateShopRequest) (*domain.Shop, error) {
	shop := &domain.Shop{
		UserID: req.UserID,
		Name:   req.Name,
	}
	err := u.shopRepo.Create(ctx, shop)
	if err != nil {
		slog.ErrorContext(ctx, "[shopUsecase] Create", "createShop", err)
		return nil, err
	}

	slog.InfoContext(ctx, "[shopUsecase] Create", "shop", shop)
	return shop, nil
}

func (u *shopUsecase) GetByUserID(ctx context.Context, userID int64) (*domain.Shop, error) {
	shop, err := u.shopRepo.GetByUserID(ctx, userID)
	if err != nil {
		slog.ErrorContext(ctx, "[shopUsecase] GetByUserID", "getShop", err)
		return nil, err
	}

	return shop, nil
}
