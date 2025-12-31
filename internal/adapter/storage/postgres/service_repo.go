package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/umutaraz/pulseguard/internal/core/domain"
)

type PostgresServiceRepository struct {
	db *pgxpool.Pool
}

func NewPostgresServiceRepository(db *pgxpool.Pool) *PostgresServiceRepository {
	return &PostgresServiceRepository{
		db: db,
	}
}

func (r *PostgresServiceRepository) Create(ctx context.Context, service *domain.Service) error {
	query := `
		INSERT INTO services (id, name, url, interval_seconds, type, thresholds, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	
	thresholdsJSON, err := json.Marshal(service.Thresholds)
	if err != nil {
		return fmt.Errorf("failed to marshal thresholds: %w", err)
	}

	_, err = r.db.Exec(ctx, query,
		service.ID,
		service.Name,
		service.URL,
		int(service.Interval.Seconds()),
		service.Type,
		thresholdsJSON,
		service.Status,
		service.CreatedAt,
		service.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	return nil
}


func (r *PostgresServiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Service, error) {
	query := `SELECT id, name, url, interval_seconds, type, thresholds, status, created_at, updated_at FROM services WHERE id = $1`

	var s domain.Service
	var intervalSeconds int
	var thresholdsJSON []byte
	var statusStr string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.Name, &s.URL, &intervalSeconds, &s.Type, &thresholdsJSON, &statusStr, &s.CreatedAt, &s.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("service not found")
		}
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	s.Interval = time.Duration(intervalSeconds) * time.Second
	s.Status = domain.ServiceStatus(statusStr)

	if err := json.Unmarshal(thresholdsJSON, &s.Thresholds); err != nil {
		return nil, fmt.Errorf("failed to unmarshal thresholds: %w", err)
	}

	return &s, nil
}

func (r *PostgresServiceRepository) GetAll(ctx context.Context) ([]*domain.Service, error) {
	query := `SELECT id, name, url, interval_seconds, type, thresholds, status, created_at, updated_at FROM services`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}
	defer rows.Close()

	var services []*domain.Service

	for rows.Next() {
		var s domain.Service
		var intervalSeconds int
		var thresholdsJSON []byte
		var statusStr string

		if err := rows.Scan(&s.ID, &s.Name, &s.URL, &intervalSeconds, &s.Type, &thresholdsJSON, &statusStr, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}

		s.Interval = time.Duration(intervalSeconds) * time.Second
		s.Status = domain.ServiceStatus(statusStr)
		json.Unmarshal(thresholdsJSON, &s.Thresholds)

		services = append(services, &s)
	}

	return services, nil
}

func (r *PostgresServiceRepository) Update(ctx context.Context, service *domain.Service) error {
	query := `
		UPDATE services 
		SET status = $1, updated_at = $2 
		WHERE id = $3
	`

	// Sadece status update ediyoruz şimdilik (Analyzer için)
	_, err := r.db.Exec(ctx, query, service.Status, service.UpdatedAt, service.ID)
	if err != nil {
		return fmt.Errorf("failed to update service: %w", err)
	}
	return nil
}

func (r *PostgresServiceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM services WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
