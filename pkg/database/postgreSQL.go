package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func OpenPostgreSQL() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping to database failed: %w", err)
	}

	err = MigratePostgreSQL(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func MigratePostgreSQL(database *sql.DB) error {
	if err := goose.Up(database, "../pkg/database/migrations"); err != nil {
		return fmt.Errorf("failed to apply migrations: %v\n", err)
	}

	return nil
}
