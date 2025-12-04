package botrepo

import (
	"database/sql"
	"fmt"
	"logging_api/internal/models"

	"github.com/lib/pq"
)

type BotRepo struct {
	db *sql.DB
}

func NewBotRepo(db *sql.DB) *BotRepo {
	return &BotRepo{
		db: db,
	}
}

func (r *BotRepo) GetBotByID(botID string) (*models.Bot, error) {
	query := `
		SELECT id, code, name, bot_type, language, description, tags, owner_id, is_active, created_at, updated_at
		FROM bots
		WHERE id = $1
	`

	var bot models.Bot
	err := r.db.QueryRow(query, botID).Scan(
		&bot.ID,
		&bot.Code,
		&bot.Name,
		&bot.BotType,
		&bot.Language,
		&bot.Description,
		pq.Array(&bot.Tags),
		&bot.OwnerID,
		&bot.IsActive,
		&bot.CreatedAt,
		&bot.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &bot, nil
}

func (r *BotRepo) GetBotByCode(code string) (*models.Bot, error) {
	query := `
		SELECT id, code, name, bot_type, language, description, tags, owner_id, is_active, created_at, updated_at
		FROM bots
		WHERE code = $1
	`

	var bot models.Bot
	err := r.db.QueryRow(query, code).Scan(
		&bot.ID,
		&bot.Code,
		&bot.Name,
		&bot.BotType,
		&bot.Language,
		&bot.Description,
		pq.Array(&bot.Tags),
		&bot.OwnerID,
		&bot.IsActive,
		&bot.CreatedAt,
		&bot.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &bot, nil
}

func (r *BotRepo) CreateBot(bot *models.Bot) (*models.Bot, error) {
	query := `
		INSERT INTO bots (code, name, bot_type, language, description, tags, owner_id, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		bot.Code,
		bot.Name,
		bot.BotType,
		bot.Language,
		bot.Description,
		pq.Array(bot.Tags),
		bot.OwnerID,
		bot.IsActive,
	).Scan(&bot.ID, &bot.CreatedAt, &bot.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return bot, nil
}

func (r *BotRepo) UpdateBot(bot *models.Bot) (*models.Bot, error) {
	query := `
		UPDATE bots
		SET name = $2, bot_type = $3, language = $4, description = $5, tags = $6, is_active = $7, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		query,
		bot.ID,
		bot.Name,
		bot.BotType,
		bot.Language,
		bot.Description,
		pq.Array(bot.Tags),
		bot.IsActive,
	).Scan(&bot.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return bot, nil
}

func (r *BotRepo) DeleteBot(botID string) error {
	query := `DELETE FROM bots WHERE id = $1`

	result, err := r.db.Exec(query, botID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("bot not found")
	}

	return nil
}

func (r *BotRepo) GetBotsByOwner(ownerID string) ([]*models.Bot, error) {
	query := `
		SELECT id, code, name, bot_type, language, description, tags, owner_id, is_active, created_at, updated_at
		FROM bots
		WHERE owner_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bots []*models.Bot
	for rows.Next() {
		var bot models.Bot
		err := rows.Scan(
			&bot.ID,
			&bot.Code,
			&bot.Name,
			&bot.BotType,
			&bot.Language,
			&bot.Description,
			pq.Array(&bot.Tags),
			&bot.OwnerID,
			&bot.IsActive,
			&bot.CreatedAt,
			&bot.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		bots = append(bots, &bot)
	}

	return bots, nil
}

func (r *BotRepo) GetAllBots() ([]*models.Bot, error) {
	query := `
		SELECT id, code, name, bot_type, language, description, tags, owner_id, is_active, created_at, updated_at
		FROM bots
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bots []*models.Bot
	for rows.Next() {
		var bot models.Bot
		err := rows.Scan(
			&bot.ID,
			&bot.Code,
			&bot.Name,
			&bot.BotType,
			&bot.Language,
			&bot.Description,
			pq.Array(&bot.Tags),
			&bot.OwnerID,
			&bot.IsActive,
			&bot.CreatedAt,
			&bot.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		bots = append(bots, &bot)
	}

	return bots, nil
}

func (r *BotRepo) GetBotCodeByID(botID string) (string, error) {
	query := `SELECT code FROM bots WHERE id = $1`

	var code string
	err := r.db.QueryRow(query, botID).Scan(&code)
	if err != nil {
		return "", err
	}

	return code, nil
}

// GetBotCodeAndNameByID получает code и name бота по его ID
func (r *BotRepo) GetBotCodeAndNameByID(botID string) (code string, name string, err error) {
	query := `SELECT code, name FROM bots WHERE id = $1`

	err = r.db.QueryRow(query, botID).Scan(&code, &name)
	if err != nil {
		return "", "", err
	}

	return code, name, nil
}
