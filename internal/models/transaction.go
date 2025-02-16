package models

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionTypeTransfer TransactionType = "transfer"
	TransactionTypePurchase TransactionType = "purchase"
)

type Transaction struct {
	ID            uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Type          TransactionType `gorm:"not null"`
	FromUsername  string          `gorm:"type:string;not null"`
	ToUsername    string          `gorm:"type:string;not null"`
	Amount        int             `gorm:"type:int;not null"`
	MerchItemName *string         `gorm:"type:string"`
	CreatedAt     time.Time       `gorm:"not null"`
}
