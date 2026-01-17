package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/senatroxx/filmix-backend/internal/http/handlers"
	"github.com/senatroxx/filmix-backend/internal/http/middleware"
)

func ShowtimeRoutes(r fiber.Router, h *handlers.Handlers) {
	r.Get("/movies/:movieId/showtimes", middleware.Protected(), h.Showtime.GetShowtimesByMovie)
	r.Get("/theaters/:theaterId/showtimes", middleware.Protected(), h.Showtime.GetShowtimesByTheater)

	showtimes := r.Group("/showtimes", middleware.Protected())
	showtimes.Get("/:id", h.Showtime.GetShowtimeByID)
}
