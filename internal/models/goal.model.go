package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Goal is a savings target not bound to any specific account.
// Contributions deduct from a user-chosen account each time via an expense transaction.
type Goal struct {
	gorm.Model
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID       uuid.UUID `gorm:"type:uuid;index"`
	Name         string    `gorm:"type:varchar(100)"`
	TargetAmount float64   `gorm:"type:decimal(12,2)"`
	SavedAmount  float64   `gorm:"type:decimal(12,2)"`
}

type GoalRequest struct {
	Name         string
	TargetAmount float64
	UserID       string
}

type GoalStatus struct {
	ID        uuid.UUID
	Name      string
	Target    float64
	Current   float64
	Remaining float64
	Percent   int
	Reached   bool
}
