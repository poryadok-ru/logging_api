package bot_handler

import (
	"net/http"

	"logging_api/internal/models"
	customerrors "logging_api/internal/utils/errors"
	"logging_api/internal/utils/validator_error_handling"

	"github.com/gin-gonic/gin"
)

type BotService interface {
	CreateBot(bot *models.Bot) (*models.Bot, error)
	GetBotByID(botID string) (*models.Bot, error)
	GetBotByCode(code string) (*models.Bot, error)
	GetAllBots() ([]*models.Bot, error)
	UpdateBot(bot *models.Bot) (*models.Bot, error)
	DeleteBot(botID string) error
}

type BotHandler struct {
	botService BotService
}

func NewBotHandler(botService BotService) *BotHandler {
	return &BotHandler{
		botService: botService,
	}
}

// CreateBot создаёт нового бота
// @Summary Создать бота
// @Description Создаёт нового бота/робота (требуется админский токен)
// @Tags bots
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateBotRequest true "Данные бота"
// @Success 201 {object} models.Bot
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/bots [post]
func (h *BotHandler) CreateBot(c *gin.Context) {
	var request CreateBotRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		if validationErrs := validator_error_handling.ValidateError(err); validationErrs != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrs.Errors()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных"})
		return
	}

	bot := &models.Bot{
		Code:        request.Code,
		Name:        request.Name,
		BotType:     request.BotType,
		Language:    request.Language,
		Description: request.Description,
		Tags:        request.Tags,
		OwnerID:     request.OwnerID,
		IsActive:    request.IsActive,
	}

	createdBot, err := h.botService.CreateBot(bot)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdBot)
}

// GetBot получает бота по ID
// @Summary Получить бота
// @Description Возвращает информацию о боте по ID (требуется админский токен)
// @Tags bots
// @Produce json
// @Security BearerAuth
// @Param bot_id path string true "ID бота (UUID)"
// @Success 200 {object} models.Bot
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/bots/{bot_id} [get]
func (h *BotHandler) GetBot(c *gin.Context) {
	botID := c.Param("bot_id")

	bot, err := h.botService.GetBotByID(botID)
	if err != nil {
		if customerrors.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bot)
}

// GetAllBots получает список всех ботов
// @Summary Получить всех ботов
// @Description Возвращает список всех ботов (требуется админский токен)
// @Tags bots
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Bot
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/bots [get]
func (h *BotHandler) GetAllBots(c *gin.Context) {
	bots, err := h.botService.GetAllBots()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bots)
}

// UpdateBot обновляет данные бота
// @Summary Обновить бота
// @Description Обновляет данные бота (требуется админский токен)
// @Tags bots
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param bot_id path string true "ID бота (UUID)"
// @Param request body UpdateBotRequest true "Обновлённые данные"
// @Success 200 {object} models.Bot
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/bots/{bot_id} [put]
func (h *BotHandler) UpdateBot(c *gin.Context) {
	botID := c.Param("bot_id")

	var request UpdateBotRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		if validationErrs := validator_error_handling.ValidateError(err); validationErrs != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrs.Errors()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных"})
		return
	}

	// Получаем существующего бота
	existingBot, err := h.botService.GetBotByID(botID)
	if err != nil {
		if customerrors.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Обновляем только переданные поля
	if request.Code != nil {
		existingBot.Code = *request.Code
	}
	if request.Name != nil {
		existingBot.Name = *request.Name
	}
	if request.BotType != nil {
		existingBot.BotType = *request.BotType
	}
	if request.Language != nil {
		existingBot.Language = *request.Language
	}
	if request.Description != nil {
		existingBot.Description = request.Description
	}
	if request.Tags != nil {
		existingBot.Tags = request.Tags
	}
	if request.OwnerID != nil {
		existingBot.OwnerID = request.OwnerID
	}
	if request.IsActive != nil {
		existingBot.IsActive = *request.IsActive
	}

	updatedBot, err := h.botService.UpdateBot(existingBot)
	if err != nil {
		if customerrors.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedBot)
}

// DeleteBot удаляет бота
// @Summary Удалить бота
// @Description Полностью удаляет бота из базы данных (требуется админский токен)
// @Tags bots
// @Produce json
// @Security BearerAuth
// @Param bot_id path string true "ID бота (UUID)"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/bots/{bot_id} [delete]
func (h *BotHandler) DeleteBot(c *gin.Context) {
	botID := c.Param("bot_id")

	err := h.botService.DeleteBot(botID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "бот удалён"})
}

