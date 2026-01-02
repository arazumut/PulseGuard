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
		INSERT INTO services (id, name, url, interval, type, thresholds, status, slack_enabled, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	
	thresholdsJSON, _ := json.Marshal(service.Thresholds)

	_, err := r.db.Exec(ctx, query,
		service.ID,
		service.Name,
		service.URL,
		service.Interval,
		service.Type,
		thresholdsJSON,
		service.Status,
		service.SlackEnabled,
		service.CreatedAt,
		service.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	return nil
}


func (r *PostgresServiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Service, error) {
	query := `
		SELECT id, name, url, interval, type, thresholds, status, slack_enabled, created_at, updated_at
		FROM services
		WHERE id = $1
	`

	var service domain.Service
	var thresholdsJSON []byte

	err := r.db.QueryRow(ctx, query, id).Scan(
		&service.ID,
		&service.Name,
		&service.URL,
		&service.Interval,
		&service.Type,
		&thresholdsJSON,
		&service.Status,
		&service.SlackEnabled,
		&service.CreatedAt,
		&service.UpdatedAt,
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

	return &service, nil
}

func (r *PostgresServiceRepository) GetAll(ctx context.Context) ([]*domain.Service, error) {
	query := `
		SELECT id, name, url, interval, type, thresholds, status, slack_enabled, created_at, updated_at
		FROM services
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query services: %w", err)
	}
	defer rows.Close()

	var services []*domain.Service
	for rows.Next() {
		var service domain.Service
		var thresholdsJSON []byte

		if err := rows.Scan(
			&service.ID, &service.Name, &service.URL, &service.Interval, &service.Type, &thresholdsJSON, &service.Status, &service.SlackEnabled, &service.CreatedAt, &service.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(thresholdsJSON, &service.Thresholds); err != nil {
			// Log error but continue? Or fail? Let's continue for now
		}

		services = append(services, &service)
	}

	return services, nil
}

func (r *PostgresServiceRepository) Update(ctx context.Context, service *domain.Service) error {
	query := `
		UPDATE services 
		SET status = $1, updated_at = $2 
		WHERE id = $3
	`

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
