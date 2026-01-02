package service

import (
	"context"
	"errors"
	"net/http"
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

func (s *MonitorService) RegisterService(ctx context.Context, name, url string, interval int, slackEnabled bool) (*domain.Service, error) {
	// Validate URL
	if _, err := http.NewRequest("GET", url, nil); err != nil {
		return nil, errors.New("invalid URL")
	}

	// Default interval if invalid
	if interval < 1 {
		interval = 60
	}

	intervalDuration := time.Duration(interval) * time.Second
	service := domain.NewService(name, url, intervalDuration, slackEnabled)

	if err := s.repo.Create(ctx, service); err != nil {
		return nil, err
	}

	s.scheduler.StartMonitorForService(service)

	return service, nil
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

func (s *MonitorService) DeleteService(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.scheduler.StopMonitorForService(id)
	return nil
}
