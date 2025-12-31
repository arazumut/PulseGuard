package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/umutaraz/pulseguard/internal/adapter/handler/http"
	"github.com/umutaraz/pulseguard/internal/adapter/handler/websocket"
	"github.com/umutaraz/pulseguard/internal/adapter/notification/slack"
	"github.com/umutaraz/pulseguard/internal/adapter/storage/postgres"
	"github.com/umutaraz/pulseguard/internal/config"
	"github.com/umutaraz/pulseguard/internal/core/domain"
	"github.com/umutaraz/pulseguard/internal/core/service"
	"github.com/umutaraz/pulseguard/internal/monitor/pinger"
	"github.com/umutaraz/pulseguard/internal/monitor/scheduler"
	"github.com/umutaraz/pulseguard/pkg/logger"
)

func main() {
	// Config & Logger init
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	logger.InitLogger(cfg.App.LogLevel)
	slog.Info("Starting PulseGuard", "env", cfg.App.Environment)

	// 3. Init Adapters (Infrastructure)
	ctx := context.Background()
	dbPool, err := postgres.NewConnection(ctx, cfg.Postgres)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Repos
	repo := postgres.NewPostgresServiceRepository(dbPool)
	metricRepo := postgres.NewPostgresMetricRepository(dbPool)

	// Init Notification Service (Slack)
	slackService := slack.NewSlackService(cfg.Notification.SlackWebhookURL)

	// Init Monitoring Engine
	httpPinger := pinger.NewHTTPPinger(5 * time.Second)
	engine := scheduler.NewMonitoringEngine(repo, httpPinger)

	if err := engine.LoadAndStart(ctx); err != nil {
		slog.Error("Failed to load services from DB", "error", err)
	}

	// Init Analyzer (The Brain)
	analyzer := service.NewAnalyzerService(repo, metricRepo, slackService)
	
	// Init WebSocket Hub (The Broadcaster)
	hub := websocket.NewHub()
	go hub.Run()

	// Wire Engine results to Analyzer AND WebSocket
	engine.SetResultHandler(func(result domain.CheckResult) {
		// 1. Analyze & Save (Async)
		go analyzer.AnalyzeResult(context.Background(), result)
		
		// 2. Broadcast to Dashboard (Real-time)
		hub.BroadcastCheckResult(result)
	})

	// Inject engine into service
	monitorService := service.NewMonitorService(repo, metricRepo, engine)
	serviceHandler := http.NewServiceHandler(monitorService)

	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		AppName:      cfg.App.Name,
	})

	http.SetupRouter(app, serviceHandler)
	
	// Add WebSocket Route
	app.Use("/ws", websocket.UpgradeMiddleware)
	app.Get("/ws", websocket.NewWebSocketHandler(hub))

	go func() {
		addr := fmt.Sprintf(":%s", cfg.Server.Port)
		if err := app.Listen(addr); err != nil {
			slog.Error("Server shutdown", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exited")
}
