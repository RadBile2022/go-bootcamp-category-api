package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
)

func InitDB(connectionString string) (*sql.DB, error) {
	// Parse config dari connection string
	config, err := pgx.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}

	// Open database menggunakan stdlib adapter
	db := stdlib.OpenDB(*config)

	// Test connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(3 * time.Hour) // opsional

	log.Println("Database connected successfully")
	return db, nil
}
