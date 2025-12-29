package http

import (
	"github.com/gofiber/fiber/v2"
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
	Name     string `json:"name"`
	URL      string `json:"url"`
	Interval int    `json:"interval"` // Seconds
}

func (h *ServiceHandler) Register(c *fiber.Ctx) error {
	var req CreateServiceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	result, err := h.svc.RegisterService(c.Context(), req.Name, req.URL, req.Interval)
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
