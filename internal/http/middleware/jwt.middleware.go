package middleware

import (
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired token")
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			// Extract user_id/role if needed here, or just let it pass.
			// The JWT middleware stores the token in c.Locals("user") by default.
			return c.Next()
		},
	})
}
