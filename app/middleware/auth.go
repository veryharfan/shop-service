package middleware

import (
	"log/slog"
	"shop-service/app/domain"
	"shop-service/app/dto"
	"shop-service/pkg"
	"shop-service/pkg/ctxutil"

	"github.com/gofiber/fiber/v2"
)

func Auth(secretKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		token, err := pkg.GetTokenFromHeaders(c.Get("Authorization"))
		if err != nil {
			slog.ErrorContext(c.Context(), "[middleware] Auth", "GetTokenFromHeaders", err)
			return c.Status(fiber.StatusUnauthorized).JSON(dto.Error(domain.ErrUnauthorized))
		}

		claims, err := pkg.ParseJwtToken(token, secretKey)
		if err != nil {
			slog.ErrorContext(c.Context(), "[middleware] Auth", "ParseJwtToken", err)
			return c.Status(fiber.StatusUnauthorized).JSON(dto.Error(domain.ErrUnauthorized))
		}

		if claims.UID == 0 {
			slog.ErrorContext(c.Context(), "[middleware] Auth", "userID", "0")
			return c.Status(fiber.StatusUnauthorized).JSON(dto.Error(domain.ErrUnauthorized))
		}

		c.Locals(ctxutil.UserIDKey, claims.UID)

		if claims.SID != nil {
			c.Locals(ctxutil.ShopIDKey, *claims.SID)
		}
		return c.Next()
	}
}
