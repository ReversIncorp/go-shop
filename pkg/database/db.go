package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"marketplace/internal/domain/entities"
	"os"
)

func InitPostgres() (*gorm.DB, error) {
	databaseConn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	db, err := gorm.Open(postgres.Open(databaseConn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&entities.User{}, &entities.Store{}, &entities.Product{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
