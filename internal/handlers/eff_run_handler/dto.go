package eff_run_handler

import (
	"logging_api/internal/models"
	"time"
)

type CreateEffRunRequest struct {
	PeriodFrom *time.Time   `json:"period_from,omitempty" example:"2024-01-01T00:00:00Z"`
	PeriodTo   *time.Time   `json:"period_to,omitempty" example:"2024-01-01T01:00:00Z"`
	Status     string       `json:"status" binding:"required,oneof=success warning error" example:"success"`
	Host       *string      `json:"host,omitempty" example:"server-01"`
	Extra      models.JSONB `json:"extra,omitempty"`
}
