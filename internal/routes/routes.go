package routes

import (
	"Avito-backend-trainee-assignment-winter-2025/config"
	"Avito-backend-trainee-assignment-winter-2025/internal/handlers"
	"Avito-backend-trainee-assignment-winter-2025/internal/routes/middleware"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func SetupRoutes(
	cfg *config.Config,
	app *fiber.App,
	merchHandler handlers.MerchHandler,
	userHandler handlers.UserHandler,
	authHandler handlers.AuthHandler,
) {
	app.Static("/swagger", "./docs")
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	api := app.Group("/api")

	api.Post("/auth", authHandler.Authenticate)

	api.Get("/merch", merchHandler.GetAllMerch)
	api.Get("/merch/:name", merchHandler.GetMerchByName)

	protected := api.Group("/", middleware.AuthMiddleware(cfg))

	protected.Get("/info", userHandler.GetInfo)
	protected.Post("/sendCoin", userHandler.SendCoins)
	protected.Get("/buy/:item", userHandler.BuyItem)
	protected.Post("/addCoin/:amount", userHandler.AddCoins)
}
