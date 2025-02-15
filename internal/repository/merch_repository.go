package repository

import (
	"Avito-backend-trainee-assignment-winter-2025/internal/models"
	"errors"
	"gorm.io/gorm"
)

type MerchRepository interface {
	GetAll() ([]models.MerchItem, error)
	GetByName(merchItemName string) (*models.MerchItem, error)
	CreateMerch(item *models.MerchItem) error
	UpdateMerch(merchItemName string, item *models.MerchItem) error
	DeleteMerch(merchItemName string) error
	InitMerchItems() error
}

type MerchRepositoryImpl struct {
	db *gorm.DB
}

func NewMerchRepository(db *gorm.DB) MerchRepository {
	return &MerchRepositoryImpl{db: db}
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
			return nil, models.ErrMerchItemNotFound
		}
		return nil, err
	}
	return &item, nil
}

func (r *MerchRepositoryImpl) CreateMerch(item *models.MerchItem) error {
	return r.db.Create(item).Error
}

func (r *MerchRepositoryImpl) UpdateMerch(merchItemName string, item *models.MerchItem) error {
	result := r.db.Model(&models.MerchItem{}).Where("name = ?", merchItemName).Updates(item)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return models.ErrMerchItemNotFound
	}

	return nil
}

func (r *MerchRepositoryImpl) DeleteMerch(merchItemName string) error {
	result := r.db.Model(&models.MerchItem{}).Where("name = ?", merchItemName).Delete(&models.MerchItem{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return models.ErrMerchItemNotFound
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
				return err
			}
		}
	}
	return nil
}
