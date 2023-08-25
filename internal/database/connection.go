package database

import (
	"errors"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Initialization initializes the database and returns the database handle.
func Initialization() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=root dbname=cinemago port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return nil, errors.New("failed to connect to database")
	}
	log.Println("Successfully connected to database.")
	return db, nil
}
