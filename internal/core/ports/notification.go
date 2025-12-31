package ports

import (
	"context"

	"github.com/umutaraz/pulseguard/internal/core/domain"
)

type NotificationService interface {
	NotifyStatusChange(ctx context.Context, service *domain.Service, oldStatus, newStatus domain.ServiceStatus) error
}
