package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/raphael-guer1n/AREA/ServiceService/internal/config"
)

func Connect(cfg config.Config) *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping DB: %v", err)
	}

	return db
}
