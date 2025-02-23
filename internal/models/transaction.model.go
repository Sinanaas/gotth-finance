package models

import (
	"time"

	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Amount          float64   `gorm:"type:decimal(12,2)"`
	Description     string    `gorm:"type:text"`
	TransactionDate time.Time `gorm:"type:date"`
	// type enum
	TransactionType constants.TransactionType `gorm:"type:int"`
	// Belongs to an account
	AccountID uuid.UUID
	Account   Account `gorm:"foreignKey:AccountID"`
	// Belongs to a user
	UserID uuid.UUID
	User   User `gorm:"foreignKey:UserID"`
	// Belongs to a category
	CategoryID uuid.UUID
	Category   Category `gorm:"foreignKey:CategoryID"`
}

type TransactionRequest struct {
	Amount      float64
	Type        int
	Description string
	CategoryID  string
	Date        string
	Account     string
	UserID      string
}
