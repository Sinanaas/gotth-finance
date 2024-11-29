package main

import (
	"fmt"
	"github.com/Sinanaas/gotth-financial-tracker/internal/initializers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"log"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("❌ Could not load environment variables")
	}

	initializers.ConnectDB(&config)
}

func main() {
	err := initializers.DB.AutoMigrate(&models.User{}, &models.Category{}, &models.Recurring{}, &models.Transaction{})
	if err != nil {
		return
	}
	fmt.Println("✅ Migration complete")
}
