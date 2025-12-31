package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/umutaraz/pulseguard/internal/core/domain"
)

type InMemoryServiceRepository struct {
	mu       sync.RWMutex
	services map[uuid.UUID]*domain.Service
}

func NewInMemoryServiceRepository() *InMemoryServiceRepository {
	return &InMemoryServiceRepository{
		services: make(map[uuid.UUID]*domain.Service),
	}
}

func (r *InMemoryServiceRepository) Create(ctx context.Context, service *domain.Service) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.services[service.ID] = service
	return nil
}

func (r *InMemoryServiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	service, ok := r.services[id]
	if !ok {
		return nil, errors.New("service not found")
	}
	return service, nil
}

func (r *InMemoryServiceRepository) GetAll(ctx context.Context) ([]*domain.Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	services := make([]*domain.Service, 0, len(r.services))
	for _, s := range r.services {
		services = append(services, s)
	}
	return services, nil
}

func (r *InMemoryServiceRepository) Update(ctx context.Context, service *domain.Service) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.services[service.ID]; !ok {
		return errors.New("service not found")
	}
	r.services[service.ID] = service
	return nil
}

func (r *InMemoryServiceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, ok := r.services[id]; !ok {
		return errors.New("service not found")
	}
	delete(r.services, id)
	return nil
}
