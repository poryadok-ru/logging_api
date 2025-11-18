package eff_run_handler

import (
	"net/http"
	"time"

	"logging_api/internal/models"
	customerrors "logging_api/internal/utils/errors"
	"logging_api/internal/utils/validator_error_handling"

	"github.com/gin-gonic/gin"
)

type EffRunService interface {
	CreateEffRun(botID string, periodFrom, periodTo *time.Time, status string, host *string, extra models.JSONB) (*models.EffRun, error)
}

type EffRunHandler struct {
	effRunService EffRunService
}

func NewEffRunHandler(effRunService EffRunService) *EffRunHandler {
	return &EffRunHandler{
		effRunService: effRunService,
	}
}

// @Summary Создать запись о запуске
// @Description Создаёт новую запись о запуске бота (требуется авторизация, только для обычных токенов с bot_id)
// @Tags eff_runs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateEffRunRequest true "Данные о запуске"
// @Success 201 {object} models.EffRun
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/eff-runs [post]
func (h *EffRunHandler) CreateEffRun(c *gin.Context) {
	var request CreateEffRunRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		if validationErrs := validator_error_handling.ValidateError(err); validationErrs != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrs.Errors()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных"})
		return
	}

	botID, exists := c.Get("bot_id")
	if !exists || botID == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "для создания записи о запуске требуется токен с привязкой к боту"})
		return
	}

	botIDStr := botID.(string)
	effRun, err := h.effRunService.CreateEffRun(botIDStr, request.PeriodFrom, request.PeriodTo, request.Status, request.Host, request.Extra)
	if err != nil {
		if customerrors.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, effRun)
}
