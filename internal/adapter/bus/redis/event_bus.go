package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"github.com/umutaraz/pulseguard/internal/config"
	"github.com/umutaraz/pulseguard/internal/core/domain"
)

const ChannelName = "pulseguard:checks"

type RedisEventBus struct {
	client *redis.Client
}

func NewRedisEventBus(cfg config.RedisConfig) (*RedisEventBus, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &RedisEventBus{client: client}, nil
}

func (r *RedisEventBus) PublishCheckResult(ctx context.Context, result domain.CheckResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}
	return r.client.Publish(ctx, ChannelName, data).Err()
}

func (r *RedisEventBus) SubscribeCheckResults(ctx context.Context) (<-chan domain.CheckResult, error) {
	sub := r.client.Subscribe(ctx, ChannelName)
	
	if _, err := sub.Receive(ctx); err != nil {
		return nil, err
	}

	ch := sub.Channel()
	outCh := make(chan domain.CheckResult)

	go func() {
		defer close(outCh)
		defer sub.Close()

		for msg := range ch {
			var result domain.CheckResult
			if err := json.Unmarshal([]byte(msg.Payload), &result); err != nil {
				slog.Error("Redis: Failed to unmarshal message", "error", err)
				continue
			}
			outCh <- result
		}
	}()

	return outCh, nil
}
