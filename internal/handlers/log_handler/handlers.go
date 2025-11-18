package log_handler

import (
	"net/http"

	"logging_api/internal/models"
	"logging_api/internal/utils/validator_error_handling"

	"github.com/gin-gonic/gin"
)

type LogService interface {
	CreateLog(botID *string, status, msg string) (*models.Log, error)
}

type LogHandler struct {
	logService LogService
}

func NewLogHandler(logService LogService) *LogHandler {
	return &LogHandler{
		logService: logService,
	}
}

// @Summary Создать лог
// @Description Создаёт новый лог от имени текущего бота (требуется авторизация)
// @Tags logs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateLogRequest true "Данные лога"
// @Success 201 {object} models.Log
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/logs [post]
func (h *LogHandler) CreateLog(c *gin.Context) {
	var request CreateLogRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		if validationErrs := validator_error_handling.ValidateError(err); validationErrs != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrs.Errors()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных"})
		return
	}

	botID, exists := c.Get("bot_id")
	if !exists {
		log, err := h.logService.CreateLog(nil, request.Status, request.Msg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, log)
		return
	}

	botIDStr := botID.(string)
	log, err := h.logService.CreateLog(&botIDStr, request.Status, request.Msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, log)
}
