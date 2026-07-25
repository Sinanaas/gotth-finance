package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Goal is a savings target backed by a real account. Progress is derived from
// that account's balance — money moves via transfers, so net worth is preserved.
type Goal struct {
	gorm.Model
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID       uuid.UUID `gorm:"type:uuid;index"`
	Name         string    `gorm:"type:varchar(100)"`
	TargetAmount float64   `gorm:"type:decimal(12,2)"`
	AccountID    uuid.UUID `gorm:"type:uuid"`
	Account      Account   `gorm:"foreignKey:AccountID"`
}

type GoalRequest struct {
	Name         string
	TargetAmount float64
	AccountID    string
	UserID       string
}

// GoalStatus is a derived, read-time view. Current = the linked account balance.
type GoalStatus struct {
	ID          uuid.UUID
	Name        string
	AccountID   uuid.UUID
	AccountName string
	Target      float64
	Current     float64
	Remaining   float64
	Percent     int
	Reached     bool
}
