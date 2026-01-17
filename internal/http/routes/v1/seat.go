package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/senatroxx/filmix-backend/internal/http/handlers"
	"github.com/senatroxx/filmix-backend/internal/http/middleware"
)

func SeatRoutes(r fiber.Router, h *handlers.Handlers) {
	r.Get("/showtimes/:showtimeId/seats", middleware.Protected(), h.Seat.GetSeatsForShowtime)
}
