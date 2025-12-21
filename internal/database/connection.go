package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(databaseURL string) error {

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	DB=db
	log.Printf("Connected to database via GORM")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
