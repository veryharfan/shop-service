package usecase

import (
	"context"
	"database/sql"
	"log/slog"
	"shop-service/app/domain"
	"shop-service/config"
)

type shopUsecase struct {
	shopRepo domain.ShopRepository
	userRepo domain.UserRepository
	cfg      *config.Config
}

func NewShopUsecase(shopRepo domain.ShopRepository, userRepo domain.UserRepository, cfg *config.Config) domain.ShopUsecase {
	return &shopUsecase{shopRepo, userRepo, cfg}
}

func (u *shopUsecase) Create(ctx context.Context, req domain.CreateShopRequest) (*domain.Shop, error) {
	_, err := u.shopRepo.GetByUserID(ctx, req.UserID)
	if err == nil || err != domain.ErrNotFound {
		slog.ErrorContext(ctx, "[shopUsecase] Create", "shop already exists", req.UserID)
		return nil, domain.ErrValidation
	}

	tx, err := u.shopRepo.BeginTransaction(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "[shopUsecase] Create", "beginTransaction", err)
		return nil, err
	}

	shop := &domain.Shop{
		UserID: req.UserID,
		Name:   req.Name,
	}

	err = u.shopRepo.WithTransaction(ctx, tx, func(ctx context.Context, tx *sql.Tx) error {
		err = u.shopRepo.Create(ctx, shop, tx)
		if err != nil {
			slog.ErrorContext(ctx, "[shopUsecase] Create", "createShop", err)
			return err
		}

		err = u.userRepo.PatchUserShop(ctx, req.UserID, domain.UserShopUpdateRequest{ShopID: shop.ID})
		if err != nil {
			slog.ErrorContext(ctx, "[shopUsecase] Create", "patchUserShop", err)
			return err
		}
		return nil
	})
	if err != nil {
		slog.ErrorContext(ctx, "[shopUsecase] Create", "withTransaction", err)
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
