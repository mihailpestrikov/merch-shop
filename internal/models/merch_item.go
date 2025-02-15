package models

import (
	"github.com/google/uuid"
)

type MerchItem struct {
	ID    uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name  string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Price int       `gorm:"type:int;not null"`
}
