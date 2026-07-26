package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Budget is one monthly spending limit for a category. Spend-to-date is never
// stored — it is derived from transactions at read time (see GetBudgetStatus).
type Budget struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID     uuid.UUID `gorm:"type:uuid;index"`
	CategoryID uuid.UUID `gorm:"type:uuid"`
	Category   Category  `gorm:"foreignKey:CategoryID"`
	Amount     float64   `gorm:"type:decimal(12,2)"`
}

type BudgetRequest struct {
	CategoryID string
	Amount     float64
	UserID     string
}

// BudgetStatus is a derived, read-time view of a budget for the current month.
type BudgetStatus struct {
	ID           uuid.UUID
	CategoryName string
	Limit        float64
	Spent        float64
	Remaining    float64
	Over         bool
	Percent      int
}
