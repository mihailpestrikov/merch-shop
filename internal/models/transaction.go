package models

import (
	"github.com/google/uuid"
	"time"
)

type TransactionType string

const (
	TransactionTypeTransfer TransactionType = "transfer"
	TransactionTypePurchase TransactionType = "purchase"
)

type Transaction struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Type        TransactionType `gorm:"not null"`
	FromUserID  uuid.UUID       `gorm:"type:uuid;not null"`
	ToUserID    uuid.UUID       `gorm:"type:uuid;not null"`
	Amount      int             `gorm:"not null"`
	MerchItemID *uuid.UUID      `gorm:"type:uuid"`
	CreatedAt   time.Time       `gorm:"not null"`
}
