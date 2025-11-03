package authservice

import (
	"database/sql"
	"fmt"
	"logging_api/internal/models"
	customerrors "logging_api/internal/utils/errors"
)

type TokenInfo struct {
	TokenID  string
	BotID    string
	OwnerID  string
	IsAdmin  bool
	IsActive bool
}

type AuthRepoInterface interface {
	CreateToken(botID *string, name string, isAdmin bool) (*models.Token, error)
	GetTokenByID(tokenID string) (*models.Token, error)
	GetTokenWithOwner(tokenID string) (token *models.Token, ownerID string, err error)
	UpdateToken(tokenID, newName string) (*models.Token, error)
	DeactivateToken(tokenID string) error
	DeleteToken(tokenID string) error
}

type BotsRepoInterface interface {
	GetBotByID(botID string) (*models.Bot, error)
}

type AuthService struct {
	authRepo AuthRepoInterface
	botsRepo BotsRepoInterface
}

func NewAuthService(authRepo AuthRepoInterface, botsRepo BotsRepoInterface) *AuthService {
	return &AuthService{
		authRepo: authRepo,
		botsRepo: botsRepo,
	}
}

func (s *AuthService) CreateToken(botID *string, tokenName string, isAdmin bool) (*models.Token, error) {
	// Если токен НЕ админский, bot_id обязателен
	if !isAdmin {
		if botID == nil || *botID == "" {
			return nil, fmt.Errorf("%w: bot_id обязателен для обычных токенов", customerrors.ErrNotFound)
		}
		// Проверяем существование бота
		_, err := s.botsRepo.GetBotByID(*botID)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("%w: бот с id %s не найден", customerrors.ErrNotFound, *botID)
			}
			return nil, fmt.Errorf("ошибка проверки бота: %w", err)
		}
	}
	// Если админский токен, bot_id может быть NULL

	token, err := s.authRepo.CreateToken(botID, tokenName, isAdmin)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания токена: %w", err)
	}

	return token, nil
}

func (s *AuthService) UpdateToken(tokenID, newName string) (*models.Token, error) {
	token, err := s.authRepo.UpdateToken(tokenID, newName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: токен не найден", customerrors.ErrNotFound)
		}
		return nil, fmt.Errorf("ошибка обновления токена: %w", err)
	}

	return token, nil
}

func (s *AuthService) DeactivateToken(tokenID string) error {
	err := s.authRepo.DeactivateToken(tokenID)
	if err != nil {
		return fmt.Errorf("ошибка деактивации токена: %w", err)
	}

	return nil
}

func (s *AuthService) DeleteToken(tokenID string) error {
	err := s.authRepo.DeleteToken(tokenID)
	if err != nil {
		return fmt.Errorf("ошибка удаления токена: %w", err)
	}

	return nil
}

func (s *AuthService) ValidateToken(tokenID string) (*TokenInfo, error) {
	token, ownerID, err := s.authRepo.GetTokenWithOwner(tokenID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: токен не найден", customerrors.ErrNotFound)
		}
		return nil, fmt.Errorf("ошибка получения токена: %w", err)
	}

	botID := ""
	if token.BotID != nil {
		botID = *token.BotID
	}

	return &TokenInfo{
		TokenID:  token.ID,
		BotID:    botID,
		OwnerID:  ownerID,
		IsAdmin:  token.IsAdmin,
		IsActive: token.IsActive,
	}, nil
}

func (s *AuthService) GetMe(tokenID string) (*models.Token, error) {
	token, err := s.authRepo.GetTokenByID(tokenID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: токен не найден", customerrors.ErrNotFound)
		}
		return nil, fmt.Errorf("ошибка получения токена: %w", err)
	}

	return token, nil
}
