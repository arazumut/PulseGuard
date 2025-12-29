package service

import (
	"context"
	"errors"
	"time"

	"github.com/umutaraz/pulseguard/internal/core/domain"
	"github.com/umutaraz/pulseguard/internal/core/ports"
)

type MonitorService struct {
	repo ports.ServiceRepository
}

func NewMonitorService(repo ports.ServiceRepository) *MonitorService {
	return &MonitorService{
		repo: repo,
	}
}

func (s *MonitorService) RegisterService(ctx context.Context, name, url string, intervalSeconds int) (*domain.Service, error) {
	if name == "" || url == "" {
		return nil, errors.New("name and url are required")
	}
	if intervalSeconds < 1 {
		intervalSeconds = 60
	}

	newService := domain.NewService(name, url, time.Duration(intervalSeconds)*time.Second)

	if err := s.repo.Create(ctx, newService); err != nil {
		return nil, err
	}

	return newService, nil
}

func (s *MonitorService) ListServices(ctx context.Context) ([]*domain.Service, error) {
	return s.repo.GetAll(ctx)
}
