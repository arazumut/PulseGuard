package pinger

import (
	"context"
	"net/http"
	"time"

	"github.com/umutaraz/pulseguard/internal/core/domain"
)

type HTTPPinger struct {
	client *http.Client
}

func NewHTTPPinger(timeout time.Duration) *HTTPPinger {
	return &HTTPPinger{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (p *HTTPPinger) Ping(ctx context.Context, service *domain.Service) domain.CheckResult {
	start := time.Now()
	
	req, err := http.NewRequestWithContext(ctx, "GET", service.URL, nil)
	if err != nil {
		return domain.CheckResult{
			ServiceID:    service.ID,
			CheckedAt:    start,
			Success:      false,
			ErrorMessage: "invalid url request creation failed: " + err.Error(),
			Latency:      0,
		}
	}

	resp, err := p.client.Do(req)
	latency := time.Since(start)

	if err != nil {
		return domain.CheckResult{
			ServiceID:    service.ID,
			CheckedAt:    start,
			Success:      false,
			ErrorMessage: err.Error(),
			Latency:      latency,
		}
	}
	defer resp.Body.Close()

	success := resp.StatusCode >= 200 && resp.StatusCode < 400
	var errMsg string
	if !success {
		errMsg = http.StatusText(resp.StatusCode)
	}

	return domain.CheckResult{
		ServiceID:    service.ID,
		CheckedAt:    start,
		StatusCode:   resp.StatusCode,
		Latency:      latency,
		Success:      success,
		ErrorMessage: errMsg,
	}
}
