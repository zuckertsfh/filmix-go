package handlers

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/senatroxx/filmix-backend/internal/database/entities"
	"github.com/senatroxx/filmix-backend/internal/http/dto"
	"github.com/senatroxx/filmix-backend/internal/services"
	"github.com/senatroxx/filmix-backend/internal/utilities"
)

type ShowtimeHandler struct {
	showtimeService services.IShowtimeService
}

func NewShowtimeHandler(showtimeService services.IShowtimeService) *ShowtimeHandler {
	return &ShowtimeHandler{showtimeService: showtimeService}
}

func (h *ShowtimeHandler) GetShowtimesByMovie(c *fiber.Ctx) error {
	movieIDParam := c.Params("movieId")
	movieID, err := uuid.Parse(movieIDParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid movie ID")
	}

	var dateFilter *time.Time
	dateParam := c.Query("date")
	if dateParam != "" {
		parsed, err := time.Parse("2006-01-02", dateParam)
		if err == nil {
			dateFilter = &parsed
		}
	}

	showtimes, err := h.showtimeService.GetShowtimesByMovieID(c.Context(), movieID, dateFilter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch showtimes")
	}

	response := h.mapShowtimesToResponse(showtimes, false)
	return utilities.NewSuccessResponse(c, http.StatusOK, "Showtimes retrieved successfully", response)
}

func (h *ShowtimeHandler) GetShowtimesByTheater(c *fiber.Ctx) error {
	theaterIDParam := c.Params("theaterId")
	theaterID, err := uuid.Parse(theaterIDParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid theater ID")
	}

	var dateFilter *time.Time
	dateParam := c.Query("date")
	if dateParam != "" {
		parsed, err := time.Parse("2006-01-02", dateParam)
		if err == nil {
			dateFilter = &parsed
		}
	}

	showtimes, err := h.showtimeService.GetShowtimesByTheaterID(c.Context(), theaterID, dateFilter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch showtimes")
	}

	response := h.mapShowtimesToResponse(showtimes, true)
	return utilities.NewSuccessResponse(c, http.StatusOK, "Showtimes retrieved successfully", response)
}

func (h *ShowtimeHandler) GetShowtimeByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid showtime ID")
	}

	showtime, err := h.showtimeService.GetShowtimeByID(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Showtime not found")
	}

	response := h.mapShowtimeToResponse(showtime, true)
	return utilities.NewSuccessResponse(c, http.StatusOK, "Showtime retrieved successfully", response)
}

func (h *ShowtimeHandler) mapShowtimesToResponse(showtimes []entities.Showtime, includeMovie bool) []dto.ShowtimeResponse {
	var response []dto.ShowtimeResponse
	for _, st := range showtimes {
		response = append(response, h.mapShowtimeToResponse(&st, includeMovie))
	}
	return response
}

func (h *ShowtimeHandler) mapShowtimeToResponse(st *entities.Showtime, includeMovie bool) dto.ShowtimeResponse {
	resp := dto.ShowtimeResponse{
		ID:        st.ID,
		Time:      st.Time,
		ExpiredAt: st.ExpiredAt,
	}

	if st.Studio != nil {
		resp.Studio = dto.StudioResponse{
			ID:   st.Studio.ID,
			Name: st.Studio.Name,
		}
	}

	if st.Theater != nil {
		resp.Theater = dto.TheaterResponse{
			ID:        st.Theater.ID,
			Name:      st.Theater.Name,
			Address:   st.Theater.Address,
			Latitude:  st.Theater.Latitude,
			Longitude: st.Theater.Longitude,
		}
		if st.Theater.Cinema != nil {
			resp.Theater.Cinema = dto.CinemaResponse{
				ID:      st.Theater.Cinema.ID,
				Name:    st.Theater.Cinema.Name,
				LogoURL: st.Theater.Cinema.LogoURL,
			}
		}
	}

	if st.Pricing != nil {
		resp.Price = st.Pricing.Price
	}

	if includeMovie && st.Movie != nil {
		resp.Movie = &dto.MovieBrief{
			ID:        st.Movie.ID,
			Title:     st.Movie.Title,
			PosterURL: st.Movie.PosterURL,
			Duration:  st.Movie.Duration,
		}
	}

	return resp
}
