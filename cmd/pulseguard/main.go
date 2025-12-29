package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/umutaraz/pulseguard/internal/adapter/handler/http"
	"github.com/umutaraz/pulseguard/internal/adapter/storage/memory"
	"github.com/umutaraz/pulseguard/internal/config"
	"github.com/umutaraz/pulseguard/internal/core/service"
	"github.com/umutaraz/pulseguard/internal/monitor/pinger"
	"github.com/umutaraz/pulseguard/internal/monitor/scheduler"
	"github.com/umutaraz/pulseguard/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	logger.InitLogger(cfg.App.LogLevel)
	slog.Info("Starting PulseGuard", "env", cfg.App.Environment)

	repo := memory.NewInMemoryServiceRepository()

	// Init Monitoring Engine
	httpPinger := pinger.NewHTTPPinger(5 * time.Second)
	engine := scheduler.NewMonitoringEngine(repo, httpPinger)

	// Inject engine into service
	monitorService := service.NewMonitorService(repo, engine)
	serviceHandler := http.NewServiceHandler(monitorService)

	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		AppName:      cfg.App.Name,
	})

	http.SetupRouter(app, serviceHandler)

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
