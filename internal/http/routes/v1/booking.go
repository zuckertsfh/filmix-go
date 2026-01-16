package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/senatroxx/filmix-backend/internal/http/handlers"
	"github.com/senatroxx/filmix-backend/internal/http/middleware"
)

func BookingRoutes(r fiber.Router, h *handlers.Handlers) {
	bookings := r.Group("/bookings", middleware.Protected())

	bookings.Post("/", h.Booking.CreateBooking)
	bookings.Get("/", h.Booking.GetUserBookings)
	bookings.Get("/:id", h.Booking.GetBooking)
}
