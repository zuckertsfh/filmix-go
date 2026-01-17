package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/senatroxx/filmix-backend/internal/http/handlers"
	"github.com/senatroxx/filmix-backend/internal/http/middleware"
)

func MovieRoutes(r fiber.Router, h *handlers.Handlers) {
	movies := r.Group("/movies", middleware.Protected())

	movies.Get("/", h.Movie.GetAllMovies)
	movies.Get("/now-playing", h.Movie.GetNowPlaying)
	movies.Get("/:id", h.Movie.GetMovieByID)
}
