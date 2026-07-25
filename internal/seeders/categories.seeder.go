package seeders

import (
	"log"

	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"gorm.io/gorm"
)

func SeedCategories(db *gorm.DB) {
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
		{
			Name:        "Transfer Out",
			Description: "Money transferred out to another account",
		},
		{
			Name:        "Transfer In",
			Description: "Money transferred in from another account",
		},
	}

	// Idempotent: only insert system categories (user_id IS NULL) that don't exist yet,
	// so new system categories reach databases seeded by older versions.
	seeded := 0
	for _, category := range categories {
		var existing models.Category
		err := db.Where("name = ? AND user_id IS NULL AND deleted_at IS NULL", category.Name).First(&existing).Error
		if err == nil {
			continue
		}
		if err != gorm.ErrRecordNotFound {
			log.Printf("❌ Failed to check category %q: %v", category.Name, err)
			continue
		}
		if err := db.Create(&category).Error; err != nil {
			log.Printf("❌ Failed to seed category %q: %v", category.Name, err)
			continue
		}
		seeded++
	}
	if seeded > 0 {
		log.Printf("Seeded %d categories", seeded)
	}
}
