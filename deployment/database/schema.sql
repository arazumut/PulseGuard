CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS services (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    url TEXT NOT NULL,
    interval_seconds INTEGER NOT NULL DEFAULT 60,
    type VARCHAR(50) NOT NULL DEFAULT 'HTTP',
    thresholds JSONB NOT NULL DEFAULT '{"latency_warning": 500000000, "latency_critical": 2000000000}',
    status VARCHAR(50) NOT NULL DEFAULT 'UNKNOWN',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS checks (
    id BIGSERIAL PRIMARY KEY,
    service_id UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    checked_at TIMESTAMP WITH TIME ZONE NOT NULL,
    status_code INTEGER,
    latency_ns BIGINT NOT NULL,
    success BOOLEAN NOT NULL,
    error_message TEXT
);

CREATE INDEX idx_checks_service_date ON checks(service_id, checked_at DESC);
