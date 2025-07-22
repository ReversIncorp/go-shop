package database

import (
	"fmt"
	"os"

	"github.com/ztrue/tracerr"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenPostgreSQL() (*gorm.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, tracerr.Errorf("failed to connect to database: %w", err)
	}

	// Не используем AutoMigrate! Миграции только через goose.

	return db, nil
}
