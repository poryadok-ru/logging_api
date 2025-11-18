package models

import "time"

// Token представляет токен аутентификации для API
type Token struct {
	ID        string    `json:"id" db:"id" example:"550e8400-e29b-41d4-a716-446655440000" swaggertype:"string" format:"uuid"`
	BotID     *string   `json:"bot_id,omitempty" db:"bot_id" example:"550e8400-e29b-41d4-a716-446655440000" swaggertype:"string" format:"uuid"`
	Name      string    `json:"name" db:"name" binding:"required,min=3,max=100" example:"Production Server"`
	IsActive  bool      `json:"is_active" db:"is_active" example:"true"`
	IsAdmin   bool      `json:"is_admin" db:"is_admin" example:"false"`
	CreatedAt time.Time `json:"created_at" db:"created_at" example:"2023-01-15T12:00:00Z"`
}
