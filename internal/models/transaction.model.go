package models

import (
	"time"

	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	// Belongs to a user
	UserID uuid.UUID
	User   User `gorm:"foreignKey:UserID"`
	// type enum
	TransactionType constants.TransactionType `gorm:"type:int"`
	Amount          float64                   `gorm:"type:decimal(12,2)"`
	Description     string
	// Belongs to a category
	CategoryID      uuid.UUID
	Category        Category  `gorm:"foreignKey:CategoryID"`
	TransactionDate time.Time `gorm:"type:date"`
}

type TransactionRequest struct {
	Amount      float64
	Type        int
	Description string
	Category    string
	Date        string
}

type TransactionWithCategory struct {
	Amount          float64
	TransactionType constants.TransactionType
	Description     string
	CategoryName    string
	TransactionDate time.Time
}
