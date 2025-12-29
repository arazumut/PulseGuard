package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRouter(app *fiber.App, handler *ServiceHandler) {
	app.Use(logger.New())
	app.Use(cors.New())

	api := app.Group("/api/v1")

	services := api.Group("/services")
	services.Post("/", handler.Register)
	services.Get("/", handler.List)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
}
