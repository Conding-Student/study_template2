package routers

import (
	"chatbot/pkg/controllers/healthchecks"
	"chatbot/pkg/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupPrivateRoutes(app *fiber.App) {

	// app.Use(handler.BasicAuth)

	// Endpoints
	apiEndpoint := app.Group("/cagabay")
	publicEndpoint := apiEndpoint.Group("/private")
	v1Endpoint := publicEndpoint.Group("/v1", handler.AuthMiddleware, handler.ServerSwitchMain)
	v1Endpoint.Post("/api/:appVersion", handler.API)
	v1Endpoint.Post("/endpoints/:deviceid", handler.APIEndPoints)

	// Service health check
	v1Endpoint.Get("/", healthchecks.CheckServiceHealth)

}
