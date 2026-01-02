package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
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
		INSERT INTO checks (id, service_id, checked_at, status_code, latency, success, error_message)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	// ... (Rest of logic is fine, latencyNs is passed as arg 4)

	newID := uuid.New()
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
		newID,
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

// GetHistory Last N metrics for a service
func (r *PostgresMetricRepository) GetHistory(ctx context.Context, serviceID uuid.UUID, limit int) ([]domain.CheckResult, error) {
	query := `
		SELECT checked_at, status_code, latency, success, error_message 
		FROM checks 
		WHERE service_id = $1 
		ORDER BY checked_at DESC 
		LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, serviceID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query history: %w", err)
	}
	defer rows.Close()

	var results []domain.CheckResult
	for rows.Next() {
		var r domain.CheckResult
		r.ServiceID = serviceID
		var errorMessage *string
		var statusCode *int
		var latencyNs int64
		
		if err := rows.Scan(&r.CheckedAt, &statusCode, &latencyNs, &r.Success, &errorMessage); err != nil {
			return nil, err
		}

		r.Latency = time.Duration(latencyNs) // Convert int64 -> Duration

		if errorMessage != nil {
			r.ErrorMessage = *errorMessage
		}
		if statusCode != nil {
			r.StatusCode = *statusCode
		}

		results = append(results, r)
	}

	return results, nil
}

func (r *PostgresMetricRepository) GetStats(ctx context.Context, serviceID uuid.UUID, since time.Time) (*domain.ServiceStats, error) {
	query := `
		SELECT 
			COUNT(*) as total, 
			COUNT(*) FILTER (WHERE success = false) as failed,
			COALESCE(AVG(latency), 0) as avg_latency
		FROM checks
		WHERE service_id = $1 AND checked_at >= $2
	`

	var total, failed int
	var avgLatency float64

	err := r.db.QueryRow(ctx, query, serviceID, since).Scan(&total, &failed, &avgLatency)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregation aggregation stats: %w", err)
	}

	stats := &domain.ServiceStats{
		TotalChecks:  total,
		FailedChecks: failed,
		AvgLatency:   time.Duration(avgLatency),
		Since:        since,
	}

	if total > 0 {
		successCount := total - failed
		stats.UptimePercentage = (float64(successCount) / float64(total)) * 100
	} else {
		stats.UptimePercentage = 100.0
	}

	return stats, nil
}
