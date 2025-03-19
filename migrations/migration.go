package migrations

import (
	"log"

	"go-restapi-crud/models"

	"gorm.io/gorm"
)

// MigrateDB runs database migrations
func MigrateDB(db *gorm.DB) {
	log.Println("Running database migrations...")

	// Auto migrate the schema
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migrations completed successfully")
}
