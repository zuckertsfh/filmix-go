package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/senatroxx/filmix-backend/internal/http/handlers"
	"github.com/senatroxx/filmix-backend/internal/http/middleware"
)

func PaymentMethodRoutes(r fiber.Router, pmHandler *handlers.PaymentMethodHandler) {
	r.Get("/payment-methods", middleware.Protected(), pmHandler.GetPaymentMethods)
}
