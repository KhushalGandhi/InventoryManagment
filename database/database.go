package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"inventory-management/models"
	"log"
	"os"
)

var DB *gorm.DB

// ConnectDB initializes the database connection and returns a Database instance
func ConnectDB() *Database {
	dsn := os.Getenv("DB_DSN")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	return &Database{
		Store: DB,
	}
}

// Database struct for storing DB connection
type Database struct {
	Store *gorm.DB
}

// RunMigrations runs the auto-migrations for your models
func (d *Database) RunMigrations() error {
	err := d.Store.AutoMigrate(
		&models.User{},
		&models.Product{},
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
		return err
	}

	log.Println("Migration completed")
	return nil
}
