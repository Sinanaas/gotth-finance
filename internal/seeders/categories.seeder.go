package seeders

import (
	"log"

	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"gorm.io/gorm"
)

func SeedCategories(db *gorm.DB) {
	var count int64
	db.Model(&models.Category{}).Count(&count)
	if count > 0 {
		return
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Fatal("❌ Transaction failed and rolled back")
		}
	}()

	categories := []models.Category{
		{
			Name:        "Healthcare",
			Description: "Expenses related to healthcare and medical services",
		},
		{
			Name:        "Education",
			Description: "Expenses related to education and learning",
		},
		{
			Name:        "Entertainment",
			Description: "Expenses related to entertainment and leisure",
		},
		{
			Name:        "Transportation",
			Description: "Expenses related to transportation and vehicle maintenance",
		},
		{
			Name:        "Food & Beverage",
			Description: "Expenses related to food and beverages",
		},
		{
			Name:        "Shopping",
			Description: "Expenses related to shopping and retail",
		},
		{
			Name:        "Utilities",
			Description: "Expenses related to utilities and household services",
		},
		{
			Name:        "Rent & Mortgage",
			Description: "Expenses related to rent and mortgage payments",
		},
		{
			Name:        "Income",
			Description: "Income",
		},
		{
			Name:        "Initial",
			Description: "Initial account balance",
		},
	}

	for _, category := range categories {
		if err := tx.Create(&category).Error; err != nil {
			tx.Rollback()
			log.Fatal("Failed to seed categories")
		}
	}
	tx.Commit()
	log.Println("Categories seeded")
}
