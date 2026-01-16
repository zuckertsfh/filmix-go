package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/senatroxx/filmix-backend/internal/http/handlers"
	v1 "github.com/senatroxx/filmix-backend/internal/http/routes/v1"
)

func SetupRoutes(r fiber.Router, h *handlers.Handlers) {
	v1api := r.Group("/v1")

	v1.AuthRoutes(v1api, h)
}
