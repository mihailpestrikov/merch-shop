package main

import (
	"Avito-backend-trainee-assignment-winter-2025/config"
	"Avito-backend-trainee-assignment-winter-2025/internal/database"
	"Avito-backend-trainee-assignment-winter-2025/internal/handlers"
	"Avito-backend-trainee-assignment-winter-2025/internal/logger"
	"Avito-backend-trainee-assignment-winter-2025/internal/repository"
	"Avito-backend-trainee-assignment-winter-2025/internal/routes"
	"Avito-backend-trainee-assignment-winter-2025/internal/service"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		// panic before loading logger
		panic(err)
	}

	log := logger.InitLogger(cfg)

	db, err := database.ConnectDB(cfg, log)
	if err != nil {
		os.Exit(1)
	}
	defer database.CloseDB(db, log)

	err = database.RunMigrations(db, log)
	if err != nil {
		os.Exit(1)
	}

	merchRepo := repository.NewMerchRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	userRepo := repository.NewUserRepository(transactionRepo, db)

	merchService := service.NewMerchService(merchRepo, log)
	userService := service.NewUserService(userRepo, transactionRepo, merchService, log)
	authService := service.NewAuthService(db, log, userService, cfg)

	authHandler := handlers.NewAuthHandler(authService, userService, log)
	userHandler := handlers.NewUserHandler(authService, userService)
	merchHandler := handlers.NewMerchHandler(merchService)

	err = merchService.InitMerchItems()
	if err != nil {
		return
	}

	app := fiber.New()

	routes.SetupRoutes(cfg, app, merchHandler, userHandler, authHandler)

	port := cfg.AppHost + ":" + cfg.AppPort
	log.Info().Msgf("Starting server on %s", port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := app.Listen(port); err != nil {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	<-quit
	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Error().Err(err).Msg("Error shutting down server")
	} else {
		log.Info().Msg("Server stopped gracefully")
	}
}
