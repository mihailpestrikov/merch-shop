package repository

import (
	"Avito-backend-trainee-assignment-winter-2025/internal/models"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type MerchRepository interface {
	GetAll() ([]models.MerchItem, error)
	GetByName(merchItemName string) (*models.MerchItem, error)
	Create(item *models.MerchItem) error
	Update(merchItemName string, item *models.MerchItem) error
	Delete(merchItemName string) error
	InitMerchItems() error
}

type MerchRepositoryImpl struct {
	db  *gorm.DB
	log *zerolog.Logger
}

func NewMerchRepository(db *gorm.DB, log *zerolog.Logger) MerchRepository {
	return &MerchRepositoryImpl{db: db, log: log}
}

func (r *MerchRepositoryImpl) GetAll() ([]models.MerchItem, error) {
	var items []models.MerchItem
	err := r.db.Find(&items).Error
	return items, err
}

func (r *MerchRepositoryImpl) GetByName(merchItemName string) (*models.MerchItem, error) {
	var item models.MerchItem
	if err := r.db.First(&item, "name = ?", merchItemName).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("merch item with name %s not found", merchItemName)
		}
		return nil, err
	}
	return &item, nil
}

func (r *MerchRepositoryImpl) Create(item *models.MerchItem) error {
	return r.db.Create(item).Error
}

func (r *MerchRepositoryImpl) Update(merchItemName string, item *models.MerchItem) error {
	result := r.db.Model(&models.MerchItem{}).Where("name = ?", merchItemName).Updates(item)

	if result.Error != nil {
		return fmt.Errorf("failed to update merch item: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("merch item with name %s not found", merchItemName)
	}

	return nil
}

func (r *MerchRepositoryImpl) Delete(merchItemName string) error {
	result := r.db.Model(&models.MerchItem{}).Where("name = ?", merchItemName).Delete(&models.MerchItem{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete merch item: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("merch item with name %s not found", merchItemName)
	}

	return nil
}

func (r *MerchRepositoryImpl) InitMerchItems() error {
	merchCatalog := []models.MerchItem{
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

	for _, item := range merchCatalog {
		var existing models.MerchItem
		err := r.db.Where("name = ?", item.Name).First(&existing).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := r.db.Create(&item).Error; err != nil {
				// log ?
				return fmt.Errorf("error creating merch item %s: %e", item.Name, err)
			}
			r.log.Info().Msgf("Item %s created", item.Name)
			if err != nil {
				return fmt.Errorf("merch item %s already creaeted", item.Name)
			}
		}
	}
	return nil
}
