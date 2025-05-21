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

type shopHandler struct {
	shopUsecase domain.ShopUsecase
	validator   *validator.Validate
	cfg         *config.Config
}

func NewShopHandler(shopUsecase domain.ShopUsecase, validator *validator.Validate, cfg *config.Config) *shopHandler {
	return &shopHandler{shopUsecase, validator, cfg}
}

func (h *shopHandler) Create(c *fiber.Ctx) error {
	var req domain.CreateShopRequest
	var err error

	if err = c.BodyParser(&req); err != nil {
		slog.ErrorContext(c.Context(), "[shopHandler] Create", "bodyParser", err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(domain.ErrBadRequest))
	}

	req.UserID, err = ctxutil.GetUserIDCtx(c.Context())
	if err != nil {
		slog.ErrorContext(c.Context(), "[shopHandler] Create", "getUserIDCtx", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Error(domain.ErrInternal))
	}

	if err := h.validator.Struct(req); err != nil {
		slog.ErrorContext(c.Context(), "[shopHandler] Create", "validation", err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(domain.ErrValidation))
	}

	shop, err := h.shopUsecase.Create(c.Context(), req)
	if err != nil {
		slog.ErrorContext(c.Context(), "[shopHandler] Create", "usecase", err)
		status, dto := dto.FromError(err)
		return c.Status(status).JSON(dto)
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Success(shop))
}

func (h *shopHandler) GetByUserID(c *fiber.Ctx) error {
	userID, err := ctxutil.GetUserIDCtx(c.Context())
	if err != nil {
		slog.ErrorContext(c.Context(), "[shopHandler] GetByUserID", "getUserIDCtx", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Error(domain.ErrInternal))
	}

	shop, err := h.shopUsecase.GetByUserID(c.Context(), userID)
	if err != nil {
		slog.ErrorContext(c.Context(), "[shopHandler] GetByUserID", "usecase", err)
		status, dto := dto.FromError(err)
		return c.Status(status).JSON(dto)
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(shop))
}
