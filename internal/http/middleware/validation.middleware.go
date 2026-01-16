package middleware

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

// Middleware factory: pass in a struct type
func ValidateBody[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body T

		// Parse request body
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON: " + err.Error(),
			})
		}

		// Run validator
		if err := validate.Struct(body); err != nil {
			errors := make(map[string]string)
			for _, e := range err.(validator.ValidationErrors) {
				errors[e.Field()] = e.Tag()
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"validation_errors": errors,
			})
		}

		// Store the parsed + validated struct in context
		c.Locals("body", body)

		return c.Next()
	}
}
