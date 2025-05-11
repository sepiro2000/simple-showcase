package database

import (
	"database/sql"
	"fmt"
	"time"

	"backend/config"

	_ "github.com/go-sql-driver/mysql"
)

// ConnectWriteDB establishes a connection to the write database
func ConnectWriteDB(cfg *config.Config) (*sql.DB, error) {
	return connectDB(cfg, cfg.WriteDBHost)
}

// ConnectReadDB establishes a connection to the read database
func ConnectReadDB(cfg *config.Config) (*sql.DB, error) {
	return connectDB(cfg, cfg.ReadDBHost)
}

// connectDB establishes a connection to the database with the given host
func connectDB(cfg *config.Config, host string) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser,
		cfg.DBPassword,
		host,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}
