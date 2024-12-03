package models

import "github.com/google/uuid"

type Account struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	// User relation
	UserID uuid.UUID `gorm:"type:uuid;not null"`
	User   User      `gorm:"foreignKey:UserID;references:ID"`
	// Transaction relation
	Transactions []Transaction `gorm:"foreignKey:AccountID"`
	// Recurring relation
	Recurrings []Recurring `gorm:"foreignKey:AccountID"`
	// Account fields
	Name        string  `gorm:"type:varchar(100);not null"`
	Description string  `gorm:"type:text"`
	Balance     float64 `gorm:"type:decimal(10,2);default:0.00"`
}
