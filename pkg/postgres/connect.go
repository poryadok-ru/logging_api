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

	// Connection pool для высокой нагрузки
	db.SetMaxOpenConns(100)           // Максимум 100 открытых соединений
	db.SetMaxIdleConns(25)            // 25 соединений в idle пуле
	db.SetConnMaxLifetime(time.Hour)  // Переоткрывать соединения каждый час
	db.SetConnMaxIdleTime(10 * time.Minute) // Закрывать простаивающие соединения через 10 минут

	return db, nil
}
