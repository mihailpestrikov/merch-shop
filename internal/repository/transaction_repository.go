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

func NewTransactionRepository(db *gorm.DB) *TransactionRepositoryImpl {
	return &TransactionRepositoryImpl{db: db}
}

func (r *TransactionRepositoryImpl) CreateTransaction(tx *gorm.DB, transaction *models.Transaction) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	if err := db.Create(transaction).Error; err != nil {
		return err
	}
	return nil
}
