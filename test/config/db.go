package config

import (
	"fmt"
	"log"
	"os"

	"github.com/sukhinho/sukhbat/test/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ .env not loaded:", err)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to PostgreSQL:", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Message{})
	if err != nil {
		log.Fatal("❌ AutoMigrate failed:", err)
	}

	DB = db
	log.Println("✅ Connected to PostgreSQL & migrated tables")
}
