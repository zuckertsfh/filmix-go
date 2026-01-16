package http

import (
	"database/sql"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/rs/zerolog"
	"github.com/senatroxx/filmix-backend/internal/config"
	"github.com/senatroxx/filmix-backend/internal/http/handlers"
	"github.com/senatroxx/filmix-backend/internal/http/middleware"
	"github.com/senatroxx/filmix-backend/internal/http/routes"
)

type API struct {
	App    *fiber.App
	Config *config.Config
	Logger zerolog.Logger
	Wg     *sync.WaitGroup
}

func InitializeAPI(cfg *config.Config, h *handlers.Handlers, db *sql.DB, log zerolog.Logger) *API {
	// This function would typically set up the API routes and handlers.
	fiberConfig := config.NewFiberConfig()
	// Override basic timeouts if needed, or keep them in NewFiberConfig
	fiberConfig.IdleTimeout = time.Minute
	fiberConfig.ReadTimeout = 10 * time.Second
	fiberConfig.WriteTimeout = 30 * time.Second
	fiberConfig.AppName = "Filmix Backend"

	app := fiber.New(fiberConfig)

	app.Use(middleware.RequestLogger(cfg.Mode, log))
	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/live",
		ReadinessProbe: func(c *fiber.Ctx) bool {
			if err := db.PingContext(c.Context()); err != nil {
				return false
			}
			return true
		},
		ReadinessEndpoint: "/ready",
	}))

	api := app.Group("/api")
	routes.SetupRoutes(api, h)

	return &API{
		App:    app,
		Config: cfg,
		Logger: log,
	}
}

func (a *API) Run() {
	go func() {
		if err := a.App.Listen(":" + a.Config.Port); err != nil {
			a.Logger.Fatal().Err(err).Msg("Failed to start server.")
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	_ = <-c // This blocks the main thread until an interrupt is received
	a.Logger.Info().Msg("Gracefully shutting down...")
	_ = a.App.Shutdown()

	a.Logger.Info().Msg("Running cleanup tasks...")

	// Your cleanup tasks go here
	// db.Close()
	// redisConn.Close()
	a.Logger.Info().Msg("Fiber was successful shutdown.")
}
