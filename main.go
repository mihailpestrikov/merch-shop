package main

import (
	"Avito-backend-trainee-assignment-winter-2025/config"
	"Avito-backend-trainee-assignment-winter-2025/internal/database"
	"Avito-backend-trainee-assignment-winter-2025/internal/logger"
	"Avito-backend-trainee-assignment-winter-2025/internal/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		/*
			If an error occurs at the initialization stage,
			there is no point in continuing execution, so panic can be used to stop
		*/
		panic(err)
	}

	log := logger.InitLogger(cfg)

	db, err := database.ConnectDB(cfg, log)
	if err != nil {
		panic(err)
	}

	err = service.InitMerchItems(db, log)
	if err != nil {
		panic(err)
	}

	//app := fiber.New()
	//
	//app.Get("/", func(c *fiber.Ctx) error {
	//	return c.JSON(fiber.Map{"message": "Hello, Fiber!"})
	//})
	//
	//port := cfg.AppHost + ":" + cfg.AppPort
	//log.Info().Msgf("Starting server on %s", port)
	//if err := app.Listen(port); err != nil {
	//	log.Fatal().Err(err).Msg("Failed to start server")
	//}
}
