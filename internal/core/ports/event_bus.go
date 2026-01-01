package ports

import (
	"context"

	"github.com/umutaraz/pulseguard/internal/core/domain"
)

type EventBus interface {
	PublishCheckResult(ctx context.Context, result domain.CheckResult) error
	SubscribeCheckResults(ctx context.Context) (<-chan domain.CheckResult, error)
}
