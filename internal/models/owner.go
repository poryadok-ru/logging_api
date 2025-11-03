package models

import "time"

// Owner представляет владельца ботов в системе
type Owner struct {
	ID        string    `json:"id" db:"id" example:"550e8400-e29b-41d4-a716-446655440000" swaggertype:"string" format:"uuid"` // Уникальный идентификатор владельца (UUID)
	FullName  string    `json:"full_name" db:"full_name" binding:"required" example:"Иван Иванов"`                            // Полное имя владельца
	IsActive  bool      `json:"is_active" db:"is_active" example:"true"`                                                      // Активен ли владелец
	CreatedAt time.Time `json:"created_at" db:"created_at" example:"2023-01-15T12:00:00Z"`                                    // Дата и время создания записи
}
