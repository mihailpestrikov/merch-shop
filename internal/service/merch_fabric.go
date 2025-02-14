package service

import (
	"Avito-backend-trainee-assignment-winter-2025/internal/models"
	"errors"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

var merchCatalog = []*models.MerchItem{
	{Name: "t-shirt", Price: 80},
	{Name: "cup", Price: 20},
	{Name: "book", Price: 50},
	{Name: "pen", Price: 10},
	{Name: "powerbank", Price: 200},
	{Name: "hoody", Price: 300},
	{Name: "umbrella", Price: 200},
	{Name: "socks", Price: 10},
	{Name: "wallet", Price: 50},
	{Name: "pink-hoody", Price: 500},
}

func InitMerchItems(db *gorm.DB, log *zerolog.Logger) error {
	for _, item := range merchCatalog {
		var existing models.MerchItem
		err := db.Where("name = ?", item.Name).First(&existing).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(item).Error; err != nil {
				log.Error().Err(err).Msgf("Error creating merch item: %s", item.Name)
				return err
			}
			log.Info().Msgf("Merch item %s added to DB", item.Name)
		}
	}
	return nil
}
