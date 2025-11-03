package middleware

import (
	"net/http"
	"strings"

	authservice "logging_api/internal/service/auth_service"
	customerrors "logging_api/internal/utils/errors"

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

func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "токен не предоставлен"})
			c.Abort()
			return
		}

		tokenInfo, err := m.tokenService.ValidateToken(token)
		if err != nil {
			if customerrors.IsNotFound(err) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "недействительный токен"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка проверки токена"})
			}
			c.Abort()
			return
		}

		if !tokenInfo.IsActive {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "токен деактивирован"})
			c.Abort()
			return
		}

		c.Set("token_id", tokenInfo.TokenID)
		c.Set("bot_id", tokenInfo.BotID)
		c.Set("is_admin", tokenInfo.IsAdmin)

		c.Next()
	}
}

func (m *AuthMiddleware) AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "токен не предоставлен"})
			c.Abort()
			return
		}

		tokenInfo, err := m.tokenService.ValidateToken(token)
		if err != nil {
			if customerrors.IsNotFound(err) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "недействительный токен"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка проверки токена"})
			}
			c.Abort()
			return
		}

		if !tokenInfo.IsActive {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "токен деактивирован"})
			c.Abort()
			return
		}

		if !tokenInfo.IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "требуются права администратора"})
			c.Abort()
			return
		}

		c.Set("token_id", tokenInfo.TokenID)
		c.Set("bot_id", tokenInfo.BotID)
		c.Set("owner_id", tokenInfo.OwnerID)
		c.Set("is_admin", tokenInfo.IsAdmin)

		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	bearerToken := c.GetHeader("Authorization")
	if bearerToken == "" {
		return ""
	}

	parts := strings.Split(bearerToken, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}

	return bearerToken
}
