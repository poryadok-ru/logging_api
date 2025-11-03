package logrepo

import (
	"database/sql"
	"fmt"
	"logging_api/internal/models"
)

type LogRepo struct {
	db *sql.DB
}

func NewLogRepo(db *sql.DB) *LogRepo {
	return &LogRepo{db: db}
}

func (r *LogRepo) CreateLog(botID *string, status, msg string) (*models.Log, error) {
	query := `
		INSERT INTO logs (bot_id, status, msg, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, bot_id, status, msg, created_at
	`

	var log models.Log
	err := r.db.QueryRow(query, botID, status, msg).Scan(
		&log.ID,
		&log.BotID,
		&log.Status,
		&log.Msg,
		&log.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create log: %w", err)
	}

	return &log, nil
}

