package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App          AppConfig          `mapstructure:"app"`
	Server       ServerConfig       `mapstructure:"server"`
	Postgres     PostgresConfig     `mapstructure:"postgres"`
	Redis        RedisConfig        `mapstructure:"redis"`
	Notification NotificationConfig `mapstructure:"notification"`
}

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment"`
	LogLevel    string `mapstructure:"log_level"`
}

type ServerConfig struct {
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type NotificationConfig struct {
	SlackWebhookURL string `mapstructure:"slack_webhook_url"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (*Config, error) {
	v := viper.New()

	// 1. Default Values
	v.SetDefault("app.name", "PulseGuard")
	v.SetDefault("app.environment", "dev")
	v.SetDefault("app.log_level", "info")
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.read_timeout", "10s")
	v.SetDefault("server.write_timeout", "10s")

	v.SetDefault("postgres.host", "localhost")
	v.SetDefault("postgres.port", "5432")
	v.SetDefault("postgres.dbname", "pulseguard")
	v.SetDefault("postgres.user", "pulseguard")
	v.SetDefault("postgres.password", "pulseguard_password")
	v.SetDefault("postgres.sslmode", "disable")

	v.SetDefault("redis.addr", "localhost:6379")

	// 2. Config File (Support local dev)
	v.AddConfigPath(".") // Current directory
	v.AddConfigPath("./configs")
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// 3. Environment Variables (Priority: High)
	// Supports PULSE_APP_NAME, PULSE_POSTGRES_HOST etc.
	v.SetEnvPrefix("PULSE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 4. Render/Cloud Specific Bindings (No prefix)
	// These override everything if present
	v.BindEnv("postgres.host", "POSTGRES_HOST")
	v.BindEnv("postgres.user", "POSTGRES_USER")
	v.BindEnv("postgres.password", "POSTGRES_PASSWORD")
	v.BindEnv("postgres.dbname", "POSTGRES_DB")
	v.BindEnv("postgres.port", "POSTGRES_PORT")
	v.BindEnv("redis.addr", "REDIS_URL") // Render often provides full REDIS_URL
	v.BindEnv("notification.slack_webhook_url", "SLACK_WEBHOOK_URL")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	return &cfg, nil
}
