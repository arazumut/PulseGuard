package domain

import (
	"time"

	"github.com/google/uuid"
)

type ServiceStatus string

const (
	StatusHealthy  ServiceStatus = "HEALTHY"
	StatusWarning  ServiceStatus = "WARNING"
	StatusCritical ServiceStatus = "CRITICAL"
	StatusDown     ServiceStatus = "DOWN"
	StatusUnknown  ServiceStatus = "UNKNOWN"
)

type ServiceThresholds struct {
	LatencyWarning  time.Duration `json:"latency_warning"`
	LatencyCritical time.Duration `json:"latency_critical"`
}

type Service struct {
	ID        uuid.UUID     `json:"id"`
	Name      string        `json:"name"`
	URL       string        `json:"url"`
	Interval  time.Duration `json:"interval"`
	Type       string            `json:"type"`
	Thresholds ServiceThresholds `json:"thresholds"`
	Status     ServiceStatus     `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type CheckResult struct {
	ServiceID   uuid.UUID     `json:"service_id"`
	CheckedAt   time.Time     `json:"checked_at"`
	StatusCode  int           `json:"status_code"`
	Latency     time.Duration `json:"latency"`
	Success     bool          `json:"success"`
	ErrorMessage string       `json:"error_message,omitempty"`
}

type ServiceStats struct {
	UptimePercentage float64       `json:"uptime_percentage"`
	AvgLatency       time.Duration `json:"avg_latency"`
	TotalChecks      int           `json:"total_checks"`
	FailedChecks     int           `json:"failed_checks"`
	Since            time.Time     `json:"since"`
}

func NewService(name, url string, interval time.Duration) *Service {
	return &Service{
		ID:        uuid.New(),
		Name:      name,
		URL:       url,
		Interval:  interval,
		Type:     "HTTP",
		Thresholds: ServiceThresholds{
			LatencyWarning:  500 * time.Millisecond,
			LatencyCritical: 2000 * time.Millisecond,
		},
		Status:    StatusUnknown,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
