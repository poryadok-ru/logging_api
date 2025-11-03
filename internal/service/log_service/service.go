package logservice

import (
	"fmt"
	"logging_api/internal/models"
)

type LogRepoInterface interface {
	CreateLog(botID *string, status, msg string) (*models.Log, error)
}

type LogService struct {
	logRepo LogRepoInterface
}

func NewLogService(logRepo LogRepoInterface) *LogService {
	return &LogService{
		logRepo: logRepo,
	}
}

func (s *LogService) CreateLog(botID *string, status, msg string) (*models.Log, error) {
	log, err := s.logRepo.CreateLog(botID, status, msg)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания лога: %w", err)
	}
	return log, nil
}

