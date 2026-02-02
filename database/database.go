package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func InitDB() (*sql.DB, error) {
	// Ambil dari env var
	connStr := os.Getenv("DB_CONN")
	if connStr == "" {
		return nil, fmt.Errorf("DB_CONN environment variable is empty")
	}

	// Jika format key-value, konversi ke format yang pgx bisa pahami
	if strings.Contains(connStr, "host=") && !strings.Contains(connStr, "postgres://") {
		// Format sudah benar untuk pgx, langsung pakai
		log.Println("Using key-value connection string format")
	} else if strings.HasPrefix(connStr, "postgres://") {
		log.Println("Using URL connection string format")
	} else {
		return nil, fmt.Errorf("invalid connection string format")
	}

	// Open connection
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection dengan timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool untuk Supabase
	db.SetMaxOpenConns(10) // Supabase free tier max 20 connections
	db.SetMaxIdleConns(5)  // Keep 5 idle connections
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	log.Println("✅ Database connected successfully to Supabase")
	return db, nil
}
