package effrunrepo

import (
	"database/sql"
	"fmt"
	"logging_api/internal/models"
	"time"
)

type EffRunRepo struct {
	db *sql.DB
}

func NewEffRunRepo(db *sql.DB) *EffRunRepo {
	return &EffRunRepo{db: db}
}

func (r *EffRunRepo) CreateEffRun(botID string, periodFrom, periodTo *time.Time, status string, host *string, extra models.JSONB) (*models.EffRun, error) {
	query := `
		INSERT INTO eff_runs (bot_id, period_from, period_to, status, host, extra, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		RETURNING id, bot_id, period_from, period_to, status, host, extra, created_at
	`

	var effRun models.EffRun
	err := r.db.QueryRow(query, botID, periodFrom, periodTo, status, host, extra).Scan(
		&effRun.ID,
		&effRun.BotID,
		&effRun.PeriodFrom,
		&effRun.PeriodTo,
		&effRun.Status,
		&effRun.Host,
		&effRun.Extra,
		&effRun.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create eff_run: %w", err)
	}

	return &effRun, nil
}

