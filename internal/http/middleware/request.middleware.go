package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func RequestLogger(env string, logger zerolog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		stop := time.Since(start)

		clientIP := c.Get("X-Forwarded-For")
		if clientIP == "" {
			clientIP = c.IP()
		}

		method := c.Method()
		url := c.OriginalURL()
		status := c.Response().StatusCode()

		// Log line
		msg := fmt.Sprintf("%s | [%s] %s %s (%d) %v",
			clientIP, method, c.Protocol(), url, status, stop)

		if env == "prod" {
			event := logger.With().Str("type", "request").Logger()
			switch {
			case status >= 500:
				event.Error().Int("status", status).Msg(msg)
			case status >= 400:
				event.Warn().Int("status", status).Msg(msg)
			default:
				event.Info().Int("status", status).Msg(msg)
			}
		} else {
			switch {
			case status >= 500:
				logger.Error().Msg(msg)
			case status >= 400:
				logger.Warn().Msg(msg)
			default:
				logger.Info().Msg(msg)
			}
		}

		return err
	}
}
