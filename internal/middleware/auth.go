package middleware

import (
	"net/http"
	"strings"
	"time"

	authservice "logging_api/internal/service/auth_service"
	customerrors "logging_api/internal/utils/errors"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

type TokenService interface {
	ValidateToken(tokenID string) (*authservice.TokenInfo, error)
}

type AuthMiddleware struct {
	tokenService TokenService
}

func NewAuthMiddleware(tokenService TokenService) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService: tokenService,
	}
}

// validateAndSetToken выполняет общую валидацию токена и устанавливает базовые поля в контекст
func (m *AuthMiddleware) validateAndSetToken(c *gin.Context) (*authservice.TokenInfo, bool) {
	token := extractToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "токен не предоставлен"})
		c.Abort()
		return nil, false
	}

	tokenInfo, err := m.tokenService.ValidateToken(token)
	if err != nil {
		if customerrors.IsNotFound(err) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "недействительный токен"})
		} else {
			sentry.CaptureException(err)
			sentry.Flush(2 * time.Second)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка проверки токена"})
		}
		c.Abort()
		return nil, false
	}

	if !tokenInfo.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "токен деактивирован"})
		c.Abort()
		return nil, false
	}

	c.Set("token_id", tokenInfo.TokenID)
	// Устанавливаем bot_id только если он не пустой (для админских токенов bot_id может быть пустым)
	if tokenInfo.BotID != "" {
		c.Set("bot_id", tokenInfo.BotID)
	}
	c.Set("is_admin", tokenInfo.IsAdmin)

	return tokenInfo, true
}

func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := m.validateAndSetToken(c); !ok {
			return
		}
		c.Next()
	}
}

func (m *AuthMiddleware) AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenInfo, ok := m.validateAndSetToken(c)
		if !ok {
			return
		}

		if !tokenInfo.IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "требуются права администратора"})
			c.Abort()
			return
		}

		c.Set("owner_id", tokenInfo.OwnerID)
		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	bearerToken := c.GetHeader("Authorization")
	if bearerToken == "" {
		return ""
	}

	parts := strings.SplitN(bearerToken, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}

	return bearerToken
}
