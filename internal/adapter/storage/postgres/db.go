package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/umutaraz/pulseguard/internal/config"
)

func NewConnection(ctx context.Context, cfg config.PostgresConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	slog.Info("Connected to PostgreSQL", "host", cfg.Host, "db", cfg.DBName)

	if err := runMigrations(ctx, pool); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return pool, nil
}

func runMigrations(ctx context.Context, db *pgxpool.Pool) error {
	slog.Info("Checking database schema...")

	queries := []string{
		`CREATE TABLE IF NOT EXISTS services (
			id UUID PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			url VARCHAR(255) NOT NULL,
			interval BIGINT NOT NULL,
			type VARCHAR(50) NOT NULL,
			thresholds JSONB NOT NULL,
			status VARCHAR(50) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL,
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS checks (
			id UUID PRIMARY KEY,
			service_id UUID REFERENCES services(id) ON DELETE CASCADE,
			status_code INT NOT NULL,
			latency BIGINT NOT NULL,
			success BOOLEAN NOT NULL,
			error_message TEXT,
			checked_at TIMESTAMP WITH TIME ZONE NOT NULL
		);`,
		`CREATE INDEX IF NOT EXISTS idx_checks_service_id_checked_at ON checks(service_id, checked_at DESC);`,
	}

	for _, q := range queries {
		if _, err := db.Exec(ctx, q); err != nil {
			return err
		}
	}

	migrationQuery := `
		DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='services' AND column_name='slack_enabled') THEN 
				ALTER TABLE services ADD COLUMN slack_enabled BOOLEAN DEFAULT TRUE; 
			END IF; 
		END $$;
	`
	if _, err := db.Exec(ctx, migrationQuery); err != nil {
		return fmt.Errorf("failed to migrate slack_enabled: %w", err)
	}

	slog.Info("Database schema verified")
	return nil
}
