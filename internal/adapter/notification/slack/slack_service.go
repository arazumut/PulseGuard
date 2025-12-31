package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/umutaraz/pulseguard/internal/core/domain"
)

type SlackService struct {
	webhookURL string
	httpClient *http.Client
}

func NewSlackService(webhookURL string) *SlackService {
	return &SlackService{
		webhookURL: webhookURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type slackMessage struct {
	Text        string       `json:"text"`
	Attachments []attachment `json:"attachments"`
}

type attachment struct {
	Color string `json:"color"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

func (s *SlackService) NotifyStatusChange(ctx context.Context, service *domain.Service, oldStatus, newStatus domain.ServiceStatus) error {
	if s.webhookURL == "" {
		// Slack not configured, skip
		return nil
	}

	color := "#36a64f" // Green
	if newStatus == domain.ServiceStatusCritical || newStatus == domain.ServiceStatusDown {
		color = "#dc3545" // Red
	} else if newStatus == domain.ServiceStatusWarning {
		color = "#ffc107" // Yellow
	}

	msg := slackMessage{
		Text: fmt.Sprintf("Service Status Changed: *%s*", service.Name),
		Attachments: []attachment{
			{
				Color: color,
				Title: fmt.Sprintf("%s -> %s", oldStatus, newStatus),
				Text:  fmt.Sprintf("Service: %s\nURL: %s\nTime: %s", service.Name, service.URL, time.Now().Format(time.RFC3339)),
			},
		},
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal slack message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send slack notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("slack API returned error status: %d", resp.StatusCode)
	}

	slog.Info("Notification sent to Slack", "service", service.Name, "status", newStatus)
	return nil
}
