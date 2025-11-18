package models

import "time"

// EffRun представляет информацию о запуске бота за определённый период
type EffRun struct {
	ID         string     `json:"id" db:"id" example:"550e8400-e29b-41d4-a716-446655440000" swaggertype:"string" format:"uuid"`
	BotID      string     `json:"bot_id" db:"bot_id" binding:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000" swaggertype:"string" format:"uuid"`
	PeriodFrom *time.Time `json:"period_from,omitempty" db:"period_from" example:"2023-01-15T10:00:00Z"`
	PeriodTo   *time.Time `json:"period_to,omitempty" db:"period_to" example:"2023-01-15T12:00:00Z"`
	Status     string     `json:"status" db:"status" binding:"oneof=success warning error" example:"success" enums:"success,warning,error"`
	Host       *string    `json:"host,omitempty" db:"host" example:"server-01"`
	Extra      JSONB      `json:"extra,omitempty" db:"extra" swaggertype:"object"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at" example:"2023-01-15T12:00:00Z"`
}
