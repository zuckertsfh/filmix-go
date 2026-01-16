package config

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/senatroxx/filmix-backend/internal/utilities"
)

func NewFiberConfig() fiber.Config {
	return fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := "Internal Server Error"

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
				message = e.Message
			}

			return c.Status(code).JSON(utilities.BaseResponse{
				Code:    code,
				Message: message,
				Data:    nil,
			})
		},
	}
}
