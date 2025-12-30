package service

import (
	"context"
	"log/slog"

	"github.com/umutaraz/pulseguard/internal/core/domain"
	"github.com/umutaraz/pulseguard/internal/core/ports"
)

type AnalyzerService struct {
	repo ports.ServiceRepository
}

func NewAnalyzerService(repo ports.ServiceRepository) *AnalyzerService {
	return &AnalyzerService{
		repo: repo,
	}
}

func (s *AnalyzerService) AnalyzeResult(ctx context.Context, result domain.CheckResult) {
	service, err := s.repo.GetByID(ctx, result.ServiceID)
	if err != nil {
		slog.Error("Analyzer: Service not found", "service_id", result.ServiceID, "error", err)
		return
	}

	newStatus := s.determineStatus(service, result)

	if service.Status != newStatus {
		oldStatus := service.Status
		service.Status = newStatus
		service.UpdatedAt = result.CheckedAt

		slog.Info("State Transition", 
			"service", service.Name, 
			"old", oldStatus, 
			"new", newStatus,
			"msg", result.ErrorMessage,
		)

		if err := s.repo.Update(ctx, service); err != nil {
			slog.Error("Failed to update service status", "id", service.ID, "error", err)
		}
	}
}

func (s *AnalyzerService) determineStatus(service *domain.Service, result domain.CheckResult) domain.ServiceStatus {
	if !result.Success {
		return domain.StatusDown
	}

	if result.StatusCode >= 500 {
		return domain.StatusDown
	}
	if result.Latency >= service.Thresholds.LatencyCritical {
		return domain.StatusCritical
	}
	if result.Latency >= service.Thresholds.LatencyWarning {
		return domain.StatusWarning
	}
	return domain.StatusHealthy
}
