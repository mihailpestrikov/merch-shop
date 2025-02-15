package service

import (
	"Avito-backend-trainee-assignment-winter-2025/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	TransferCoins(fromUserID, toUserID uuid.UUID, amount int) error
	GetBalance(userID uuid.UUID) error
	PurchaseMerch(userID uuid.UUID, merchItemID uuid.UUID) error
	GetPurchasedItems(userID uuid.UUID) ([]models.MerchItem, error)
	GetTransactionHistory(userID uuid.UUID) ([]models.Transaction, error)
}

type UserServiceImpl struct {
	db *gorm.DB
}

func NewUserServiceImpl(db *gorm.DB) UserService {
	return &UserServiceImpl{db: db}
}

func (s *UserServiceImpl) GetBalance(userID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s *UserServiceImpl) TransferCoins(fromUserID, toUserID uuid.UUID, amount int) error {
	//TODO implement me
	panic("implement me")
}

func (s *UserServiceImpl) PurchaseMerch(userID uuid.UUID, merchItemID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s *UserServiceImpl) GetPurchasedItems(userID uuid.UUID) ([]models.MerchItem, error) {
	var items []models.MerchItem
	if err := s.db.Joins("JOIN transactions ON transactions.merch_item_id = merch_items.id").
		Where("transactions.to_user_id = ? AND transactions.type = ?", userID, models.TransactionTypePurchase).
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (s *UserServiceImpl) GetTransactionHistory(userID uuid.UUID) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := s.db.Where("from_user_id = ? OR to_user_id = ?", userID, userID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
