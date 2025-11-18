package auth_handler

import (
	"net/http"

	"logging_api/internal/models"
	customerrors "logging_api/internal/utils/errors"
	"logging_api/internal/utils/validator_error_handling"

	"github.com/gin-gonic/gin"
)

type AuthService interface {
	CreateToken(botID *string, tokenName string, isAdmin bool) (*models.Token, error)
	UpdateToken(tokenID, newName string) (*models.Token, error)
	DeactivateToken(tokenID string) error
	DeleteToken(tokenID string) error
	GetMe(tokenID string) (*models.Token, error)
}

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// @Summary Создать новый токен
// @Description Создаёт новый токен аутентификации для указанного бота (требуется админский токен)
// @Tags tokens
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateTokenRequest true "Данные для создания токена"
// @Success 201 {object} TokenResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/tokens [post]
func (h *AuthHandler) CreateToken(c *gin.Context) {
	var request CreateTokenRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		if validationErrs := validator_error_handling.ValidateError(err); validationErrs != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": validationErrs.Errors(),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "неверный формат данных",
		})
		return
	}

	token, err := h.authService.CreateToken(request.BotID, request.TokenName, request.IsAdmin)
	if err != nil {
		if customerrors.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, TokenResponse{Token: token.ID})
}

// @Summary Обновить токен
// @Description Обновляет название существующего токена (требуется админский токен)
// @Tags tokens
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param token_id path string true "ID токена (UUID)"
// @Param request body UpdateTokenRequest true "Новое название токена"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/tokens/{token_id} [put]
func (h *AuthHandler) UpdateToken(c *gin.Context) {
	tokenID := c.Param("token_id")

	var request UpdateTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		if validationErrs := validator_error_handling.ValidateError(err); validationErrs != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": validationErrs.Errors(),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "неверный формат данных",
		})
		return
	}

	token, err := h.authService.UpdateToken(tokenID, request.TokenName)
	if err != nil {
		if customerrors.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "токен обновлён",
		"token":   token,
	})
}

// @Summary Деактивировать токен
// @Description Деактивирует токен (мягкое удаление, токен перестаёт работать, требуется админский токен)
// @Tags tokens
// @Produce json
// @Security BearerAuth
// @Param token_id path string true "ID токена (UUID)"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/tokens/{token_id}/deactivate [patch]
func (h *AuthHandler) DeactivateToken(c *gin.Context) {
	tokenID := c.Param("token_id")

	err := h.authService.DeactivateToken(tokenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "токен деактивирован",
	})
}

// @Summary Удалить токен
// @Description Полностью удаляет токен из базы данных (требуется админский токен)
// @Tags tokens
// @Produce json
// @Security BearerAuth
// @Param token_id path string true "ID токена (UUID)"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/tokens/{token_id} [delete]
func (h *AuthHandler) DeleteToken(c *gin.Context) {
	tokenID := c.Param("token_id")

	err := h.authService.DeleteToken(tokenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "токен удалён",
	})
}

// @Summary Получить информацию о токене
// @Description Возвращает информацию о токене из заголовка Authorization
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Token
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	tokenID, exists := c.Get("token_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "токен не найден в контексте"})
		return
	}

	token, err := h.authService.GetMe(tokenID.(string))
	if err != nil {
		if customerrors.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, token)
}
