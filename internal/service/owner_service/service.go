package ownerservice

import (
	"database/sql"
	"fmt"
	"logging_api/internal/models"
	customerrors "logging_api/internal/utils/errors"
)

type OwnerRepoInterface interface {
	CreateOwner(fullName string, isActive bool) (*models.Owner, error)
	GetOwnerByID(ownerID string) (*models.Owner, error)
	GetAllOwners() ([]*models.Owner, error)
	UpdateOwner(ownerID string, fullName *string, isActive *bool) (*models.Owner, error)
	DeleteOwner(ownerID string) error
}

type OwnerService struct {
	ownerRepo OwnerRepoInterface
}

func NewOwnerService(ownerRepo OwnerRepoInterface) *OwnerService {
	return &OwnerService{
		ownerRepo: ownerRepo,
	}
}

func (s *OwnerService) CreateOwner(fullName string, isActive bool) (*models.Owner, error) {
	owner, err := s.ownerRepo.CreateOwner(fullName, isActive)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания владельца: %w", err)
	}
	return owner, nil
}

func (s *OwnerService) GetOwnerByID(ownerID string) (*models.Owner, error) {
	owner, err := s.ownerRepo.GetOwnerByID(ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: владелец не найден", customerrors.ErrNotFound)
		}
		return nil, fmt.Errorf("ошибка получения владельца: %w", err)
	}
	return owner, nil
}

func (s *OwnerService) GetAllOwners() ([]*models.Owner, error) {
	owners, err := s.ownerRepo.GetAllOwners()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка владельцев: %w", err)
	}
	return owners, nil
}

func (s *OwnerService) UpdateOwner(ownerID string, fullName *string, isActive *bool) (*models.Owner, error) {
	owner, err := s.ownerRepo.UpdateOwner(ownerID, fullName, isActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: владелец не найден", customerrors.ErrNotFound)
		}
		return nil, fmt.Errorf("ошибка обновления владельца: %w", err)
	}
	return owner, nil
}

func (s *OwnerService) DeleteOwner(ownerID string) error {
	err := s.ownerRepo.DeleteOwner(ownerID)
	if err != nil {
		return fmt.Errorf("ошибка удаления владельца: %w", err)
	}
	return nil
}
