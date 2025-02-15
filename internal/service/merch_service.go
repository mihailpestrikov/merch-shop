package service

import (
	"Avito-backend-trainee-assignment-winter-2025/internal/models"
	"Avito-backend-trainee-assignment-winter-2025/internal/repository"
	"github.com/rs/zerolog"
)

type MerchService interface {
	GetAllMerchItems() ([]models.MerchItem, error)
	GetMerchItemByName(merchItemName string) (*models.MerchItem, error)
	AddMerchItem(item models.MerchItem) error
	UpdateMerchItem(merchItemName string, item models.MerchItem) error
	DeleteMerchItem(merchItemName string) error
	InitMerchItems() error
}

type MerchServiceImpl struct {
	repo repository.MerchRepository
	log  *zerolog.Logger
}

func NewMerchService(repo repository.MerchRepository, log *zerolog.Logger) MerchService {
	return &MerchServiceImpl{repo: repo, log: log}
}

func (s *MerchServiceImpl) GetAllMerchItems() ([]models.MerchItem, error) {
	items, err := s.repo.GetAll()
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to get all merch items")
		return nil, err
	}
	return items, nil
}

func (s *MerchServiceImpl) GetMerchItemByName(merchItemName string) (*models.MerchItem, error) {
	item, err := s.repo.GetByName(merchItemName)
	if err != nil {
		s.log.Error().Err(err).Msgf("Failed to get merch item by name: %s", merchItemName)
		return nil, err
	}
	return item, nil
}

func (s *MerchServiceImpl) AddMerchItem(item models.MerchItem) error {
	err := s.repo.CreateMerch(&item)
	if err != nil {
		s.log.Error().Err(err).Msgf("Failed to add merch item: %v", item)
		return err
	}
	s.log.Info().Msgf("Merch item added: %v", item)
	return nil
}

func (s *MerchServiceImpl) UpdateMerchItem(merchItemName string, item models.MerchItem) error {
	err := s.repo.UpdateMerch(merchItemName, &item)
	if err != nil {
		s.log.Error().Err(err).Msgf("Failed to update merch item with name: %s", merchItemName)
		return err
	}
	s.log.Info().Msgf("Merch item updated: %v", item)
	return nil
}

func (s *MerchServiceImpl) DeleteMerchItem(merchItemName string) error {
	err := s.repo.DeleteMerch(merchItemName)
	if err != nil {
		s.log.Error().Err(err).Msgf("Failed to delete merch item with name: %s", merchItemName)
		return err
	}
	s.log.Info().Msgf("Merch item deleted with name: %s", merchItemName)
	return nil
}

func (s *MerchServiceImpl) InitMerchItems() error {
	err := s.repo.InitMerchItems()
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to initialize merch items")
		return err
	}
	s.log.Info().Msg("Merch items initialized")
	return nil
}
