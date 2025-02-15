package database

import (
	"Avito-backend-trainee-assignment-winter-2025/internal/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB, log *zerolog.Logger) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, getMigrations())
	if err := m.Migrate(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run migrations")
		return err
	}
	log.Info().Msg("Migrations completed successfully")
	return nil
}

func getMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "001_create_merch_table",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.MerchItem{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("merch_items")
			},
		},
		{
			ID: "002_create_transaction_table",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.Transaction{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("transactions")
			},
		},
		{
			ID: "003_create_user_table",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.User{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("users")
			},
		},
	}
}
