package models

import "time"

// Bot представляет бота или автоматизированную систему
type Bot struct {
	ID          string    `json:"id" db:"id" example:"550e8400-e29b-41d4-a716-446655440000" swaggertype:"string" format:"uuid"`
	Code        string    `json:"code" db:"code" binding:"required" example:"BOT_001"`
	Name        string    `json:"name" db:"name" binding:"required" example:"Telegram Bot"`
	BotType     string    `json:"bot_type" db:"bot_type" binding:"required,oneof=AI Backend Frontend Robot" example:"Backend"`
	Language    string    `json:"language" db:"language" binding:"required,oneof=Python Go N8N PIX JS C Other" example:"Python"`
	Description *string   `json:"description,omitempty" db:"description" example:"Бот для обработки сообщений"`
	Tags        []string  `json:"tags,omitempty" db:"tags" example:"telegram,automation"`
	OwnerID     *string   `json:"owner_id,omitempty" db:"owner_id" example:"550e8400-e29b-41d4-a716-446655440000" swaggertype:"string"`
	IsActive    bool      `json:"is_active" db:"is_active" example:"true"`
	CreatedAt   time.Time `json:"created_at" db:"created_at" example:"2023-01-15T12:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at" example:"2023-01-15T12:00:00Z"`
}
