package models

import "time"

// Token представляет токен аутентификации для API
type Token struct {
	ID        string    `json:"id" db:"id" example:"550e8400-e29b-41d4-a716-446655440000" swaggertype:"string" format:"uuid"`                   // Уникальный идентификатор токена (используется как API ключ)
	BotID     *string   `json:"bot_id,omitempty" db:"bot_id" example:"550e8400-e29b-41d4-a716-446655440000" swaggertype:"string" format:"uuid"` // ID бота, для которого создан токен (NULL для админских токенов)
	Name      string    `json:"name" db:"name" binding:"required,min=3,max=100" example:"Production Server"`                                    // Название токена
	IsActive  bool      `json:"is_active" db:"is_active" example:"true"`                                                                        // Активен ли токен
	IsAdmin   bool      `json:"is_admin" db:"is_admin" example:"false"`                                                                         // Обладает ли токен административными правами
	CreatedAt time.Time `json:"created_at" db:"created_at" example:"2023-01-15T12:00:00Z"`                                                      // Дата и время создания токена
}
