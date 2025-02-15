package repository

import (
	"Avito-backend-trainee-assignment-winter-2025/internal/models"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(userID uuid.UUID) (*models.User, error)
	UpdateUserBalance(userID uuid.UUID, newBalance int) error
	GetPurchasedItems(userID uuid.UUID) ([]models.MerchItem, error)
	GetTransactionHistory(userID uuid.UUID) ([]models.Transaction, error)
}

type UserRepositoryImpl struct {
	log *zerolog.Logger
	db  *gorm.DB
}

func NewUserRepositoryImpl(log *zerolog.Logger, db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{log: log, db: db}
}

func (r *UserRepositoryImpl) GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with id: %s not found", userID)
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) UpdateUserBalance(userID uuid.UUID, newBalance int) error {
	result := r.db.Model(&models.User{}).Where("id = ?", userID).Update("balance", newBalance)

	if result.Error != nil {
		return fmt.Errorf("failed to update balance: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user with ID %s not found", userID.String())
	}

	return nil
}

func (r *UserRepositoryImpl) GetPurchasedItems(userID uuid.UUID) ([]models.MerchItem, error) {
	var items []models.MerchItem

	err := r.db.Joins("JOIN transactions ON transactions.merch_item_id = merch_items.id").
		Where("transactions.to_user_id = ? AND transactions.type = ?", userID, models.TransactionTypePurchase).
		Find(&items).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get purchased items for user %s: %w", userID.String(), err)
	}

	return items, nil
}

func (r *UserRepositoryImpl) GetTransactionHistory(userID uuid.UUID) ([]models.Transaction, error) {
	var transactions []models.Transaction

	err := r.db.Where("from_user_id = ? OR to_user_id = ?", userID, userID).
		Order("created_at DESC").
		Find(&transactions).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get transaction history for user %s: %w", userID.String(), err)
	}

	return transactions, nil
}
