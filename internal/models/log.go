package models

import "time"

// Log представляет запись лога от бота
type Log struct {
	ID        int64     `json:"id" db:"id" example:"12345"`                                                                                     // Уникальный идентификатор лога
	BotID     *string   `json:"bot_id,omitempty" db:"bot_id" example:"550e8400-e29b-41d4-a716-446655440000" swaggertype:"string" format:"uuid"` // ID бота, отправившего лог
	Status    string    `json:"status" db:"status" binding:"oneof=success warning error" example:"success" enums:"success,warning,error"`       // Статус выполнения операции
	Msg       string    `json:"msg" db:"msg" binding:"required" example:"Операция выполнена успешно"`                                           // Текст сообщения лога
	CreatedAt time.Time `json:"created_at" db:"created_at" example:"2023-01-15T12:00:00Z"`                                                      // Дата и время создания лога
}
