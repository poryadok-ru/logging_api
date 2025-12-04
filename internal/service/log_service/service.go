package logservice

import (
	"fmt"
	"log"
	"logging_api/internal/models"
	"logging_api/pkg/sentry"
)

type LogRepoInterface interface {
	CreateLog(botID *string, status, msg string) (*models.Log, error)
}

type BotRepoInterface interface {
	GetBotCodeAndNameByID(botID string) (code string, name string, err error)
}

type LogService struct {
	logRepo LogRepoInterface
	botRepo BotRepoInterface
}

func NewLogService(logRepo LogRepoInterface, botRepo BotRepoInterface) *LogService {
	return &LogService{
		logRepo: logRepo,
		botRepo: botRepo,
	}
}

func (s *LogService) CreateLog(botID *string, status, msg string) (*models.Log, error) {
	logEntry, err := s.logRepo.CreateLog(botID, status, msg)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания лога: %w", err)
	}

	if botID != nil && *botID != "" && (status == "Error" || status == "Critical") {
		projectCode, botName, err := s.botRepo.GetBotCodeAndNameByID(*botID)
		if err != nil {
			log.Printf("Не удалось получить project_code и bot_name для bot_id %s: %v", *botID, err)
			projectCode = "unknown"
			botName = "unknown"
		}
		sentry.SendLog(*botID, projectCode, botName, status, msg, logEntry.CreatedAt)
	}

	return logEntry, nil
}
