package postgres

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func Connect(host string, port int, user string, password string, dbname string, sslmode string) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(time.Hour)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db, nil
}
