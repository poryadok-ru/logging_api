package ownerrepo

import (
	"database/sql"
	"fmt"
	"logging_api/internal/models"
)

type OwnerRepo struct {
	db *sql.DB
}

func NewOwnerRepo(db *sql.DB) *OwnerRepo {
	return &OwnerRepo{db: db}
}

func (r *OwnerRepo) CreateOwner(fullName string, isActive bool) (*models.Owner, error) {
	query := `
		INSERT INTO owners (full_name, is_active, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, full_name, is_active, created_at
	`

	var owner models.Owner
	err := r.db.QueryRow(query, fullName, isActive).Scan(
		&owner.ID,
		&owner.FullName,
		&owner.IsActive,
		&owner.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create owner: %w", err)
	}

	return &owner, nil
}

func (r *OwnerRepo) GetOwnerByID(ownerID string) (*models.Owner, error) {
	query := `
		SELECT id, full_name, is_active, created_at
		FROM owners
		WHERE id = $1
	`

	var owner models.Owner
	err := r.db.QueryRow(query, ownerID).Scan(
		&owner.ID,
		&owner.FullName,
		&owner.IsActive,
		&owner.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &owner, nil
}

func (r *OwnerRepo) GetAllOwners() ([]*models.Owner, error) {
	query := `
		SELECT id, full_name, is_active, created_at
		FROM owners
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get owners: %w", err)
	}
	defer rows.Close()

	var owners []*models.Owner
	for rows.Next() {
		var owner models.Owner
		err := rows.Scan(
			&owner.ID,
			&owner.FullName,
			&owner.IsActive,
			&owner.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan owner: %w", err)
		}
		owners = append(owners, &owner)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	return owners, nil
}

func (r *OwnerRepo) UpdateOwner(ownerID string, fullName *string, isActive *bool) (*models.Owner, error) {
	query := `
		UPDATE owners
		SET 
			full_name = COALESCE($2, full_name),
			is_active = COALESCE($3, is_active)
		WHERE id = $1
		RETURNING id, full_name, is_active, created_at
	`

	var owner models.Owner
	err := r.db.QueryRow(query, ownerID, fullName, isActive).Scan(
		&owner.ID,
		&owner.FullName,
		&owner.IsActive,
		&owner.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &owner, nil
}

func (r *OwnerRepo) DeleteOwner(ownerID string) error {
	query := `DELETE FROM owners WHERE id = $1`

	result, err := r.db.Exec(query, ownerID)
	if err != nil {
		return fmt.Errorf("failed to delete owner: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("owner not found")
	}

	return nil
}
