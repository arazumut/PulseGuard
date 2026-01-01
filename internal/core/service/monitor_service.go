package service

import (
	"context"
	"errors"
	"time"
	
	"github.com/google/uuid"
	"github.com/umutaraz/pulseguard/internal/core/domain"
	"github.com/umutaraz/pulseguard/internal/core/ports"
)

type Scheduler interface {
	StartMonitorForService(service *domain.Service)
	StopMonitorForService(id uuid.UUID)
}

type MonitorService struct {
	repo       ports.ServiceRepository
	metricRepo ports.MetricRepository
	scheduler  Scheduler
}

func NewMonitorService(repo ports.ServiceRepository, metricRepo ports.MetricRepository, scheduler Scheduler) *MonitorService {
	return &MonitorService{
		repo:       repo,
		metricRepo: metricRepo,
		scheduler:  scheduler,
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

	s.scheduler.StartMonitorForService(newService)

	return newService, nil
}

func (s *MonitorService) ListServices(ctx context.Context) ([]*domain.Service, error) {
	return s.repo.GetAll(ctx)
}

func (s *MonitorService) GetServiceMetrics(ctx context.Context, serviceID uuid.UUID) ([]domain.CheckResult, error) {
	return s.metricRepo.GetHistory(ctx, serviceID, 50)
}

func (s *MonitorService) GetServiceStats(ctx context.Context, serviceID uuid.UUID) (*domain.ServiceStats, error) {
	since := time.Now().Add(-24 * time.Hour)
	return s.metricRepo.GetStats(ctx, serviceID, since)
}
