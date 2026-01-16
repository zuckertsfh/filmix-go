package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/senatroxx/filmix-backend/internal/http/dto"
	"github.com/senatroxx/filmix-backend/internal/services"
	"github.com/senatroxx/filmix-backend/internal/utilities"
)

type MovieHandler struct {
	movieService services.IMovieService
}

func NewMovieHandler(movieService services.IMovieService) *MovieHandler {
	return &MovieHandler{movieService: movieService}
}

func (h *MovieHandler) GetAllMovies(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	movies, total, err := h.movieService.GetAllMovies(c.Context(), page, limit)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch movies")
	}

	var response []dto.MovieResponse
	for _, movie := range movies {
		resp := dto.MovieResponse{
			ID:          movie.ID,
			Title:       movie.Title,
			Tagline:     movie.Tagline,
			Overview:    movie.Overview,
			PosterURL:   movie.PosterURL,
			BackdropURL: movie.BackdropURL,
			TrailerURL:  movie.TrailerURL,
			Duration:    movie.Duration,
			Popularity:  movie.Popularity,
		}
		if movie.Status != nil {
			resp.Status = movie.Status.Status
		}
		if movie.Rating != nil {
			resp.Rating = movie.Rating.Rating
		}
		response = append(response, resp)
	}

	return utilities.NewPaginatedResponse(c, http.StatusOK, "Movies retrieved successfully", response, page, limit, total)
}

func (h *MovieHandler) GetNowPlaying(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	movies, total, err := h.movieService.GetNowPlayingMovies(c.Context(), page, limit)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch movies")
	}

	var response []dto.MovieResponse
	for _, movie := range movies {
		resp := dto.MovieResponse{
			ID:          movie.ID,
			Title:       movie.Title,
			Tagline:     movie.Tagline,
			Overview:    movie.Overview,
			PosterURL:   movie.PosterURL,
			BackdropURL: movie.BackdropURL,
			TrailerURL:  movie.TrailerURL,
			Duration:    movie.Duration,
			Popularity:  movie.Popularity,
		}
		if movie.Status != nil {
			resp.Status = movie.Status.Status
		}
		if movie.Rating != nil {
			resp.Rating = movie.Rating.Rating
		}
		response = append(response, resp)
	}

	return utilities.NewPaginatedResponse(c, http.StatusOK, "Now playing movies retrieved successfully", response, page, limit, total)
}

func (h *MovieHandler) GetMovieByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid movie ID")
	}

	movie, err := h.movieService.GetMovieByID(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Movie not found")
	}

	response := dto.MovieResponse{
		ID:          movie.ID,
		Title:       movie.Title,
		Tagline:     movie.Tagline,
		Overview:    movie.Overview,
		PosterURL:   movie.PosterURL,
		BackdropURL: movie.BackdropURL,
		TrailerURL:  movie.TrailerURL,
		Duration:    movie.Duration,
		Popularity:  movie.Popularity,
	}
	if movie.Status != nil {
		response.Status = movie.Status.Status
	}
	if movie.Rating != nil {
		response.Rating = movie.Rating.Rating
	}

	for _, genre := range movie.Genres {
		response.Genres = append(response.Genres, dto.GenreResponse{
			ID:    genre.ID,
			Genre: genre.Genre,
		})
	}

	return utilities.NewSuccessResponse(c, http.StatusOK, "Movie retrieved successfully", response)
}
