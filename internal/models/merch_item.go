package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MerchItem struct {
	ID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name  string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Price int       `gorm:"type:int;not null"`
}

func (m *MerchItem) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}
