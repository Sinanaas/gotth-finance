package models

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Recurring struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string    `gorm:"type:varchar(100)"`
	Amount    float64   `gorm:"type:decimal(12,2)"`
	StartDate time.Time `gorm:"type:date"`
	// type enum
	TransactionType constants.TransactionType `gorm:"type:int"`
	// periodicity enum
	Periodicity constants.Periodicity `gorm:"type:int"`
	// Belongs to a user
	UserID uuid.UUID
	User   User `gorm:"foreignKey:UserID"`
	// Belongs to a category
	CategoryID uuid.UUID
	Category   Category `gorm:"foreignKey:CategoryID"`
	// Belongs to an account
	AccountID uuid.UUID
	Account   Account `gorm:"foreignKey:AccountID"`
}

type RecurringWithCategoryName struct {
	Amount          float64
	TransactionType constants.TransactionType
	Periodicity     constants.Periodicity
	Name            string
	CategoryName    string
	AccountName     string
}

type RecurringRequest struct {
	Name            string
	Amount          float64
	TransactionType int
	Periodicity     int
	CategoryID      string
	AccountID       string
	UserID          string
}
