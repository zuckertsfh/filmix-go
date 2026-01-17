package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/senatroxx/filmix-backend/internal/http/handlers"
	"github.com/senatroxx/filmix-backend/internal/http/middleware"
)

func AuthRoutes(r fiber.Router, h *handlers.Handlers) {
	auth := r.Group("/auth")

	auth.Post("/register", h.Auth.Register)
	auth.Post("/login", h.Auth.Login)
	auth.Post("/refresh", h.Auth.RefreshToken)

	auth.Get("/me", middleware.Protected(), h.Auth.GetProfile)
}
