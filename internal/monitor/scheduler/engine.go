package scheduler

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/umutaraz/pulseguard/internal/core/domain"
	"github.com/umutaraz/pulseguard/internal/core/ports"
	"github.com/umutaraz/pulseguard/internal/monitor/pinger"
)

type ResultHandler func(result domain.CheckResult)

type MonitoringEngine struct {
	serviceRepo    ports.ServiceRepository
	pinger         *pinger.HTTPPinger
	activeMonitors map[uuid.UUID]context.CancelFunc
	mu             sync.Mutex
	onResult       ResultHandler
}

func (e *MonitoringEngine) LoadAndStart(ctx context.Context) error {
	services, err := e.serviceRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	for _, s := range services {
		e.StartMonitorForService(s)
	}
	slog.Info("Bootstrapped monitoring engine", "count", len(services))
	return nil
}

func NewMonitoringEngine(repo ports.ServiceRepository, pinger *pinger.HTTPPinger) *MonitoringEngine {
	return &MonitoringEngine{
		serviceRepo:    repo,
		pinger:         pinger,
		activeMonitors: make(map[uuid.UUID]context.CancelFunc),

	}
}

func (e *MonitoringEngine) SetResultHandler(handler ResultHandler) {
	e.onResult = handler
}
func (e *MonitoringEngine) StartMonitorForService(service *domain.Service) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if cancel, exists := e.activeMonitors[service.ID]; exists {
		cancel()
	}

	ctx, cancel := context.WithCancel(context.Background())
	e.activeMonitors[service.ID] = cancel

	slog.Info("Started monitoring for service", "service_id", service.ID, "url", service.URL, "interval", service.Interval)

	go e.monitorLoop(ctx, service)
}

func (e *MonitoringEngine) StopMonitorForService(id uuid.UUID) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if cancel, exists := e.activeMonitors[id]; exists {
		cancel()
		delete(e.activeMonitors, id)
		slog.Info("Stopped monitoring for service", "service_id", id)
	}
}

func (e *MonitoringEngine) monitorLoop(ctx context.Context, service *domain.Service) {
	ticker := time.NewTicker(service.Interval)
	defer ticker.Stop()

	e.performCheck(ctx, service)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			e.performCheck(ctx, service)
		}
	}
}

func (e *MonitoringEngine) performCheck(ctx context.Context, service *domain.Service) {
	checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second) 
	defer cancel()

	result := e.pinger.Ping(checkCtx, service)
	
	slog.Info("Health Check", 
		"service", service.Name, 
		"url", service.URL, 
		"status_code", result.StatusCode, 
		"latency", result.Latency, 
		"success", result.Success,
	)

	if e.onResult != nil {
		e.onResult(result)
	}
}
