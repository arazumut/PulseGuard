package service

import (
	"context"
	"log/slog"

	"github.com/umutaraz/pulseguard/internal/core/domain"
	"github.com/umutaraz/pulseguard/internal/core/ports"
)

type AnalyzerService struct {
	repo       ports.ServiceRepository
	metricRepo ports.MetricRepository
	notifier   ports.NotificationService
}

func NewAnalyzerService(repo ports.ServiceRepository, metricRepo ports.MetricRepository, notifier ports.NotificationService) *AnalyzerService {
	return &AnalyzerService{
		repo:       repo,
		metricRepo: metricRepo,
		notifier:   notifier,
	}
}

func (s *AnalyzerService) AnalyzeResult(ctx context.Context, result domain.CheckResult) {
	if err := s.metricRepo.Save(ctx, &result); err != nil {
		slog.Error("Analyzer: Failed to save metric", "service_id", result.ServiceID, "error", err)
	}

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
		} else {
			// Notify success update
			if service.SlackEnabled {
				go func() {
					if err := s.notifier.NotifyStatusChange(ctx, service, oldStatus, newStatus); err != nil {
						slog.Error("Failed to send notification", "error", err, "service", service.Name)
					}
				}()
			}
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
