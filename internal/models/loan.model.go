package models

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Loan struct {
	gorm.Model
	ID              uuid.UUID                 `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TransactionType constants.TransactionType `gorm:"type:int"`
	Status          bool                      `gorm:"type:bool"`
	Description     string                    `gorm:"type:text"`
	LoanDate        time.Time                 `gorm:"type:date"`
	ToWhom          string                    `gorm:"type:text"`
	Amount          float64                   `gorm:"type:decimal(12,2)"`
	// Belongs to a user
	UserID uuid.UUID
	User   User `gorm:"foreignKey:UserID"`
	// Belongs to an account
	AccountID uuid.UUID
	Account   Account `gorm:"foreignKey:AccountID"`
	// Belongs to a category
	CategoryID uuid.UUID
	Category   Category `gorm:"foreignKey:CategoryID"`
}

type LoanCategoryAccount struct {
	ID              string
	Amount          float64
	ToWhom          string
	Description     string
	LoanDate        time.Time
	Status          bool
	TransactionType constants.TransactionType
	AccountName     string
	CategoryName    string
}

type LoanRequest struct {
	Amount          float64
	ToWhom          string
	Description     string
	LoanDate        string
	Status          bool
	TransactionType int
	CategoryID      string
	AccountID       string
	UserID          string
}
