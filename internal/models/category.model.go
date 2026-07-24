package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string     `gorm:"type:varchar(100)"`
	Description string     `gorm:"type:text"`
	UserID      *uuid.UUID `gorm:"type:uuid;index"`
}

type CategoryWithTotal struct {
	Name  string
	Total float64
}
