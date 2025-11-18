package authrepo

import (
	"database/sql"
	"fmt"
	"logging_api/internal/models"
)

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (r *AuthRepo) CreateToken(botID *string, name string, isAdmin bool) (*models.Token, error) {
	query := `
		INSERT INTO tokens (bot_id, name, is_active, is_admin)
		VALUES ($1, $2, true, $3)
		RETURNING id, bot_id, name, is_active, is_admin, created_at
	`

	var token models.Token
	err := r.db.QueryRow(query, botID, name, isAdmin).Scan(
		&token.ID,
		&token.BotID,
		&token.Name,
		&token.IsActive,
		&token.IsAdmin,
		&token.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *AuthRepo) GetTokenByID(tokenID string) (*models.Token, error) {
	query := `
		SELECT id, bot_id, name, is_active, is_admin, created_at
		FROM tokens
		WHERE id = $1
	`

	var token models.Token
	err := r.db.QueryRow(query, tokenID).Scan(
		&token.ID,
		&token.BotID,
		&token.Name,
		&token.IsActive,
		&token.IsAdmin,
		&token.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *AuthRepo) GetTokenWithOwner(tokenID string) (token *models.Token, ownerID string, err error) {
	query := `
		SELECT t.id, t.bot_id, t.name, t.is_active, t.is_admin, t.created_at, b.owner_id
		FROM tokens t
		LEFT JOIN bots b ON t.bot_id = b.id
		WHERE t.id = $1
	`

	var token2 models.Token
	var ownerIDPtr *string
	err = r.db.QueryRow(query, tokenID).Scan(
		&token2.ID,
		&token2.BotID,
		&token2.Name,
		&token2.IsActive,
		&token2.IsAdmin,
		&token2.CreatedAt,
		&ownerIDPtr,
	)
	if err != nil {
		return nil, "", err
	}

	if ownerIDPtr != nil {
		ownerID = *ownerIDPtr
	}

	return &token2, ownerID, nil
}

func (r *AuthRepo) UpdateToken(tokenID, newName string) (*models.Token, error) {
	query := `
		UPDATE tokens
		SET name = $2
		WHERE id = $1
		RETURNING id, bot_id, name, is_active, is_admin, created_at
	`

	var token models.Token
	err := r.db.QueryRow(query, tokenID, newName).Scan(
		&token.ID,
		&token.BotID,
		&token.Name,
		&token.IsActive,
		&token.IsAdmin,
		&token.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *AuthRepo) DeactivateToken(tokenID string) error {
	query := `
		UPDATE tokens
		SET is_active = false
		WHERE id = $1
	`

	result, err := r.db.Exec(query, tokenID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("token not found")
	}

	return nil
}

func (r *AuthRepo) DeleteToken(tokenID string) error {
	query := `DELETE FROM tokens WHERE id = $1`

	result, err := r.db.Exec(query, tokenID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("token not found")
	}

	return nil
}
