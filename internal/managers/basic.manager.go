package managers

import "gorm.io/gorm"

type BasicManager struct {
	DB *gorm.DB
}

func NewBasicManager(DB *gorm.DB) BasicManager {
	return BasicManager{DB}
}
