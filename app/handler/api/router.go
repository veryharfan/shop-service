package handler

import (
	"shop-service/app/middleware"
	"shop-service/config"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App, shopHandler *shopHandler, cfg *config.Config) {

	shopAPIGroup := app.Group("/shop-service").Use(middleware.Auth(cfg.Jwt.SecretKey))

	shopAPIGroup.Post("/shops", shopHandler.Create)
	shopAPIGroup.Get("/shops", shopHandler.GetByUserID)

}
