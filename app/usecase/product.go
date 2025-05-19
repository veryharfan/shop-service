package usecase

import (
	"context"
	"log/slog"
	"shop-service/app/domain"
	"shop-service/config"
)

type productUsecase struct {
	productRepo domain.ProductRepository
	stockRepo   domain.StockRepository
	cfg         *config.Config
}

func NewProductUsecase(productRepo domain.ProductRepository, stockRepo domain.StockRepository, cfg *config.Config) domain.ProductUsecase {
	return &productUsecase{productRepo, stockRepo, cfg}
}

func (u *productUsecase) CreateProduct(ctx context.Context, shopID int64, req domain.ShopProductCreateRequest) (domain.ProductCreateResponse, error) {
	// Create product
	productReq := domain.ProductCreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
		ShopID:      shopID,
	}
	productResp, err := u.productRepo.CreateProduct(ctx, productReq)
	if err != nil {
		slog.ErrorContext(ctx, "[productUsecase] CreateProduct", "createProduct", err)
		return domain.ProductCreateResponse{}, err
	}

	err = u.stockRepo.InitStock(ctx, domain.InitStockRequest{
		ShopID:    shopID,
		ProductID: productResp.ID,
	})
	if err != nil {
		slog.ErrorContext(ctx, "[productUsecase] CreateProduct", "initStock", err)
		// If stock initialization fails, we should delete the product or handle it accordingly.
		// This is a placeholder for the actual deletion logic.
		// _, deleteErr := u.productRepo.DeleteProduct(ctx, productResp.ID)
		// if deleteErr != nil {
		// 	slog.ErrorContext(ctx, "[productUsecase] CreateProduct", "deleteProduct", deleteErr)
		// 	return domain.ProductCreateResponse{}, deleteErr
		// }
		// Return the error to the caller
		// so they can handle it as needed.

		// other options could be to return the error or log it
		// but not delete the product immediately,
		// but will delete the product if product fetch fails
		return domain.ProductCreateResponse{}, err
	}

	return productResp, nil
}
