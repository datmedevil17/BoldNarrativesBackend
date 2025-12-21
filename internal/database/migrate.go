package database

import (
	"log"

	"github.com/datmedevil17/BoldNarrativesBackend/internal/models"
)

func Migrate() error {
	err := DB.AutoMigrate(&models.User{}, &models.Blog{}, &models.Comment{}, &models.Vote{}, &models.Follows{})
	if err != nil {
		log.Fatal("❌ Migration failed:", err)
	}
	log.Println("Migrations completed successfully")
	return nil
}
