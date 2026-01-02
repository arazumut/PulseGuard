package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/umutaraz/pulseguard/internal/core/service"
)

type ServiceHandler struct {
	svc *service.MonitorService
}

func NewServiceHandler(svc *service.MonitorService) *ServiceHandler {
	return &ServiceHandler{
		svc: svc,
	}
}

type CreateServiceRequest struct {
	Name         string `json:"name"`
	URL          string `json:"url"`
	Interval     int    `json:"interval"` // Seconds
	SlackEnabled bool   `json:"slack_enabled"`
}

func (h *ServiceHandler) Register(c *fiber.Ctx) error {
	var req CreateServiceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	result, err := h.svc.RegisterService(c.Context(), req.Name, req.URL, req.Interval, req.SlackEnabled)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

func (h *ServiceHandler) List(c *fiber.Ctx) error {
	services, err := h.svc.ListServices(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": services,
		"count": len(services),
	})
}

func (h *ServiceHandler) GetMetrics(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid service id"})
	}

	metrics, err := h.svc.GetServiceMetrics(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	stats, err := h.svc.GetServiceStats(c.Context(), id)
	if err != nil {
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"service_id": id,
		"history":    metrics,
		"stats":      stats,
	})
}

func (h *ServiceHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid service id"})
	}

	if err := h.svc.DeleteService(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
