package repository

import (
	"Avito-backend-trainee-assignment-winter-2025/internal/models"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(tx *gorm.DB, transaction *models.Transaction) error
}

type TransactionRepositoryImpl struct {
	db *gorm.DB
}

func NewTransactionRepositoryImpl(db *gorm.DB) *TransactionRepositoryImpl {
	return &TransactionRepositoryImpl{db: db}
}

func (r *TransactionRepositoryImpl) CreateTransaction(tx *gorm.DB, transaction *models.Transaction) error {
	if err := tx.Create(transaction).Error; err != nil {
		return err
	}
	return nil
}
