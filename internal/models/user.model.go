package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username string    `gorm:"type:varchar(100);unique_index"`
	Password string    `gorm:"type:varchar(100)"`
	Email    string    `gorm:"type:varchar(100);unique_index"`
	PhotoURL string    `gorm:"type:varchar(255)"`

	// Transactions relation
	Transactions []Transaction `gorm:"foreignKey:UserID"`
	// Recurring relation
	Recurrings []Recurring `gorm:"foreignKey:UserID"`
	// Accounts relation
	Accounts []Account `gorm:"foreignKey:UserID"`
}

type SignUpInput struct {
	Email           string
	Username        string
	Password        string
	ConfirmPassword string
}

type EditUserInput struct {
	Username string
	Email    string
	PhotoURL string
}

type SignInInput struct {
	Email    string
	Password string
}

type UserResponse struct {
	ID        uuid.UUID
	Email     string
	Username  string
	Provider  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
