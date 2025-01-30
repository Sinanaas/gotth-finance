package seeders

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/initializers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"log"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables")
	}
	initializers.ConnectDB(&config)
}
func main() {
	tx := initializers.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Fatal("? Transaction failed and rolled back")
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
	}

	for _, category := range categories {
		if err := tx.Create(&category).Error; err != nil {
			tx.Rollback()
			log.Fatal("? Failed to seed categories")
		}
	}

	tx.Commit()
}
