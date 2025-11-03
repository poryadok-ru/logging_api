package models

import "time"

// EffRun представляет информацию о запуске бота за определённый период
type EffRun struct {
	ID         string     `json:"id" db:"id" example:"550e8400-e29b-41d4-a716-446655440000" swaggertype:"string" format:"uuid"`                                 // Уникальный идентификатор запуска
	BotID      string     `json:"bot_id" db:"bot_id" binding:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000" swaggertype:"string" format:"uuid"` // ID бота
	PeriodFrom *time.Time `json:"period_from,omitempty" db:"period_from" example:"2023-01-15T10:00:00Z"`                                                        // Начало периода выполнения
	PeriodTo   *time.Time `json:"period_to,omitempty" db:"period_to" example:"2023-01-15T12:00:00Z"`                                                            // Конец периода выполнения
	Status     string     `json:"status" db:"status" binding:"oneof=success warning error" example:"success" enums:"success,warning,error"`                     // Статус выполнения запуска
	Host       *string    `json:"host,omitempty" db:"host" example:"server-01"`                                                                                 // Хост, на котором выполнялся бот
	Extra      JSONB      `json:"extra,omitempty" db:"extra" swaggertype:"object"`                                                                              // Дополнительные метаданные в формате JSON
	CreatedAt  time.Time  `json:"created_at" db:"created_at" example:"2023-01-15T12:00:00Z"`                                                                    // Дата и время создания записи
}
