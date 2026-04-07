package service

import (
	"fmt"

	"github.com/google/uuid"
	"tenant-service/internal/model"
	"tenant-service/internal/repository"
	"tenant-service/internal/utils"
)

type TenantService struct {
	repo *repository.TenantRepo
}

func NewTenantService(r *repository.TenantRepo) *TenantService {
	return &TenantService{repo: r}
}

func (s *TenantService) Register(name, email string) (*model.Tenant, error) {
	exists, err := s.repo.ExistsByEmail(email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("tenant already exists")
	}

	t := model.Tenant{
		ID:     uuid.New().String(),
		Name:   name,
		Email:  email,
		APIKey: utils.GenerateAPIKey(),
	}

	if err := s.repo.Create(t); err != nil {
		return nil, err
	}

	return &t, nil
}

func (s *TenantService) Validate(apiKey string) (*model.Tenant, error) {
	t, err := s.repo.GetByAPIKey(apiKey)
	if err != nil {
		return nil, fmt.Errorf("invalid api key")
	}
	return t, nil
}