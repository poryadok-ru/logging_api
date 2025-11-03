package botservice

import (
	"database/sql"
	"fmt"
	"logging_api/internal/models"
	customerrors "logging_api/internal/utils/errors"
)

type BotRepoInterface interface {
	CreateBot(bot *models.Bot) (*models.Bot, error)
	GetBotByID(botID string) (*models.Bot, error)
	GetBotByCode(code string) (*models.Bot, error)
	GetBotsByOwner(ownerID string) ([]*models.Bot, error)
	GetAllBots() ([]*models.Bot, error)
	UpdateBot(bot *models.Bot) (*models.Bot, error)
	DeleteBot(botID string) error
}

type BotService struct {
	botRepo BotRepoInterface
}

func NewBotService(botRepo BotRepoInterface) *BotService {
	return &BotService{
		botRepo: botRepo,
	}
}

func (s *BotService) CreateBot(bot *models.Bot) (*models.Bot, error) {
	createdBot, err := s.botRepo.CreateBot(bot)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания бота: %w", err)
	}
	return createdBot, nil
}

func (s *BotService) GetBotByID(botID string) (*models.Bot, error) {
	bot, err := s.botRepo.GetBotByID(botID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: бот не найден", customerrors.ErrNotFound)
		}
		return nil, fmt.Errorf("ошибка получения бота: %w", err)
	}
	return bot, nil
}

func (s *BotService) GetBotByCode(code string) (*models.Bot, error) {
	bot, err := s.botRepo.GetBotByCode(code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: бот не найден", customerrors.ErrNotFound)
		}
		return nil, fmt.Errorf("ошибка получения бота: %w", err)
	}
	return bot, nil
}

func (s *BotService) GetAllBots() ([]*models.Bot, error) {
	bots, err := s.botRepo.GetAllBots()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка ботов: %w", err)
	}
	return bots, nil
}

func (s *BotService) UpdateBot(bot *models.Bot) (*models.Bot, error) {
	updatedBot, err := s.botRepo.UpdateBot(bot)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: бот не найден", customerrors.ErrNotFound)
		}
		return nil, fmt.Errorf("ошибка обновления бота: %w", err)
	}
	return updatedBot, nil
}

func (s *BotService) DeleteBot(botID string) error {
	err := s.botRepo.DeleteBot(botID)
	if err != nil {
		return fmt.Errorf("ошибка удаления бота: %w", err)
	}
	return nil
}

