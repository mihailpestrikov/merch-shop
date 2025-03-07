package repository

import (
	"Avito-backend-trainee-assignment-winter-2025/internal/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(userID uuid.UUID) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	UpdateUserBalance(tx *gorm.DB, fromUsername string, newBalance int) error
	GetPurchasedItems(username string) ([]models.MerchItem, error)
	GetTransactionHistory(username string) ([]models.Transaction, error)
	CreateUser(user *models.User) error
	UpdateUser(userID uuid.UUID, user *models.User) error
	DeleteUser(userID uuid.UUID) error
}

type UserRepositoryImpl struct {
	transactionRepo TransactionRepository
	db              *gorm.DB
}

func NewUserRepository(transactionRepo TransactionRepository, db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{transactionRepo: transactionRepo, db: db}
}

func (r *UserRepositoryImpl) GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserRepositoryImpl) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserRepositoryImpl) UpdateUserBalance(tx *gorm.DB, username string, newBalance int) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	result := db.Model(&models.User{}).Where("username = ?", username).Update("balance", newBalance)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return models.ErrUserNotFound
	}

	return nil
}

func (r *UserRepositoryImpl) GetPurchasedItems(username string) ([]models.MerchItem, error) {
	var items []models.MerchItem

	err := r.db.Joins("JOIN transactions ON transactions.merch_item_name = merch_items.name").
		Where("transactions.to_username = ? AND transactions.type = ?", username, models.TransactionTypePurchase).
		Find(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *UserRepositoryImpl) GetTransactionHistory(username string) ([]models.Transaction, error) {
	var transactions []models.Transaction

	err := r.db.Where("from_username = ? OR to_username = ?", username, username).
		Order("created_at DESC").
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *UserRepositoryImpl) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepositoryImpl) UpdateUser(userID uuid.UUID, user *models.User) error {
	result := r.db.Model(&models.User{}).Where("id = ?", userID).Updates(user)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return models.ErrMerchItemNotFound
	}

	return nil
}

func (r *UserRepositoryImpl) DeleteUser(userID uuid.UUID) error {
	result := r.db.Model(&models.User{}).Where("id = ?", userID).Delete(&models.User{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return models.ErrMerchItemNotFound
	}

	return nil
}
