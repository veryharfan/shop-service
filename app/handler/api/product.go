package handler

import (
	"log/slog"
	"shop-service/app/domain"
	"shop-service/app/dto"
	"shop-service/config"
	"shop-service/pkg/ctxutil"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type productHandler struct {
	productUsecase domain.ProductUsecase
	validator      *validator.Validate
	cfg            *config.Config
}

func NewProductHandler(productUsecase domain.ProductUsecase, validator *validator.Validate, cfg *config.Config) *productHandler {
	return &productHandler{productUsecase, validator, cfg}
}

func (h *productHandler) Create(c *fiber.Ctx) error {
	var req domain.ShopProductCreateRequest
	var err error

	if err = c.BodyParser(&req); err != nil {
		slog.ErrorContext(c.Context(), "[productHandler] Create", "bodyParser", err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(domain.ErrBadRequest))
	}

	if err := h.validator.Struct(req); err != nil {
		slog.ErrorContext(c.Context(), "[productHandler] Create", "validation", err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(domain.ErrValidation))
	}

	shopID, err := ctxutil.GetShopIDCtx(c.Context())
	if err != nil {
		slog.ErrorContext(c.Context(), "[productHandler] Create", "getShopIDCtx", err)
		return c.Status(fiber.StatusUnauthorized).JSON(dto.Error(domain.ErrUnauthorized))
	}

	product, err := h.productUsecase.CreateProduct(c.Context(), shopID, req)
	if err != nil {
		slog.ErrorContext(c.Context(), "[productHandler] Create", "usecase", err)
		status, dto := dto.FromError(err)
		return c.Status(status).JSON(dto)
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Success(product))
}
