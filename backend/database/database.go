package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"sensory-navigator/config"
)

var DB *sql.DB

func Connect(cfg *config.DBConfig) error {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func GetDB() *sql.DB {
	return DB
}

