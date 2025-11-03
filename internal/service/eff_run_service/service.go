package effrunservice

import (
	"fmt"
	"logging_api/internal/models"
	"time"
)

type EffRunRepoInterface interface {
	CreateEffRun(botID string, periodFrom, periodTo *time.Time, status string, host *string, extra models.JSONB) (*models.EffRun, error)
}

type EffRunService struct {
	effRunRepo EffRunRepoInterface
}

func NewEffRunService(effRunRepo EffRunRepoInterface) *EffRunService {
	return &EffRunService{
		effRunRepo: effRunRepo,
	}
}

func (s *EffRunService) CreateEffRun(botID string, periodFrom, periodTo *time.Time, status string, host *string, extra models.JSONB) (*models.EffRun, error) {
	effRun, err := s.effRunRepo.CreateEffRun(botID, periodFrom, periodTo, status, host, extra)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания записи о запуске: %w", err)
	}
	return effRun, nil
}

