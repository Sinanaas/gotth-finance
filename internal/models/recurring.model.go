package models

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Recurring struct {
	gorm.Model
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	// Belongs to a user
	UserID      uuid.UUID
	User        User                  `gorm:"foreignKey:UserID"`
	Name        string                `gorm:"type:varchar(100)"`
	Description string                `gorm:"type:text"`
	Amount      float64               `gorm:"type:decimal(12,2)"`
	Periodicity constants.Periodicity `gorm:"type:int"`
	// Belongs to a category
	CategoryID uuid.UUID
	Category   Category `gorm:"foreignKey:CategoryID"`
	// Belongs to an account
	AccountID uuid.UUID
	Account   Account `gorm:"foreignKey:AccountID"`
}
