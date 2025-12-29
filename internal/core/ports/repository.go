package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/umutaraz/pulseguard/internal/core/domain"
)

type ServiceRepository interface {
	Create(ctx context.Context, service *domain.Service) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Service, error)
	GetAll(ctx context.Context) ([]*domain.Service, error)
	Update(ctx context.Context, service *domain.Service) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type MetricRepository interface {
	Save(ctx context.Context, result *domain.CheckResult) error
}
