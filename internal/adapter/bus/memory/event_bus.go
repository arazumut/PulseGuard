package memory

import (
	"context"
	"sync"

	"github.com/umutaraz/pulseguard/internal/core/domain"
)

type MemoryEventBus struct {
	mu          sync.RWMutex
	subscribers []chan domain.CheckResult
}

func NewMemoryEventBus() *MemoryEventBus {
	return &MemoryEventBus{
		subscribers: make([]chan domain.CheckResult, 0),
	}
}

func (m *MemoryEventBus) PublishCheckResult(ctx context.Context, result domain.CheckResult) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, ch := range m.subscribers {
		// Non-blocking send to avoid stalling if one subscriber is slow
		select {
		case ch <- result:
		default:
			// Channel full, drop message or log warning
		}
	}
	return nil
}

func (m *MemoryEventBus) SubscribeCheckResults(ctx context.Context) (<-chan domain.CheckResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	ch := make(chan domain.CheckResult, 100) // Buffer to prevent blocking
	m.subscribers = append(m.subscribers, ch)

	return ch, nil
}
