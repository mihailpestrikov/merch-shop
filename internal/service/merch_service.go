package service

import (
	"Avito-backend-trainee-assignment-winter-2025/internal/models"
	"errors"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type MerchService interface {
	GetAllMerchItems() ([]models.MerchItem, error)
	GetMerchItemByID(merchItemID uuid.UUID) (*models.MerchItem, error)
	AddMerchItem(item models.MerchItem) error
	UpdateMerchItem(merchItemID uuid.UUID, item models.MerchItem) error
	DeleteMerchItem(merchItemID uuid.UUID) error
	InitMerchItems() error
}

type MerchServiceImpl struct {
	db  *gorm.DB
	log *zerolog.Logger
}

func NewMerchServiceImpl(db *gorm.DB, log *zerolog.Logger) MerchService {
	return &MerchServiceImpl{db: db, log: log}
}

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

func (s *MerchServiceImpl) GetAllMerchItems() ([]models.MerchItem, error) {
	//TODO implement me
	panic("implement me")
}

func (s *MerchServiceImpl) GetMerchItemByID(merchItemID uuid.UUID) (*models.MerchItem, error) {
	//TODO implement me
	panic("implement me")
}

func (s *MerchServiceImpl) AddMerchItem(item models.MerchItem) error {
	//TODO implement me
	panic("implement me")
}

func (s *MerchServiceImpl) UpdateMerchItem(merchItemID uuid.UUID, item models.MerchItem) error {
	//TODO implement me
	panic("implement me")
}

func (s *MerchServiceImpl) DeleteMerchItem(merchItemID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s *MerchServiceImpl) InitMerchItems() error {
	for _, item := range merchCatalog {
		var existing models.MerchItem
		err := s.db.Where("name = ?", item.Name).First(&existing).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := s.db.Create(item).Error; err != nil {
				s.log.Error().Err(err).Msgf("Error creating merch item: %s", item.Name)
				return err
			}
			s.log.Info().Msgf("Merch item %s added to DB", item.Name)
		}
	}
	return nil
}
