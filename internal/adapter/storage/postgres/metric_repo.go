package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/umutaraz/pulseguard/internal/core/domain"
)

type PostgresMetricRepository struct {
	db *pgxpool.Pool
}

func NewPostgresMetricRepository(db *pgxpool.Pool) *PostgresMetricRepository {
	return &PostgresMetricRepository{
		db: db,
	}
}

func (r *PostgresMetricRepository) Save(ctx context.Context, result *domain.CheckResult) error {
	query := `
		INSERT INTO checks (service_id, checked_at, status_code, latency_ns, success, error_message)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	latencyNs := result.Latency.Nanoseconds()
	
	var errorMessage *string
	if result.ErrorMessage != "" {
		errorMessage = &result.ErrorMessage
	}
	var statusCode *int
	if result.StatusCode != 0 {
		statusCode = &result.StatusCode
	}

	_, err := r.db.Exec(ctx, query,
		result.ServiceID,
		result.CheckedAt,
		statusCode,
		latencyNs,
		result.Success,
		errorMessage,
	)

	if err != nil {
		return fmt.Errorf("failed to save check result: %w", err)
	}

	return nil
}

func (r *PostgresMetricRepository) GetHistory(ctx context.Context, serviceID string, limit int) ([]domain.CheckResult, error) {
	return nil, nil 
}
