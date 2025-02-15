package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username     string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	Balance      int       `gorm:"not null;default:1000"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
