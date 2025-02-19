package database

import (
	"Avito-backend-trainee-assignment-winter-2025/config"
	"fmt"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(cfg *config.Config, log *zerolog.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSslMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to database")
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	log.Info().Msg("Connected to database")

	return db, nil
}

func CloseDB(db *gorm.DB, log *zerolog.Logger) {
	dbSQL, err := db.DB()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get DB instance for closing")

		return
	}
	if err := dbSQL.Close(); err != nil {
		log.Error().Err(err).Msg("Failed to close database connection")
	}
}
