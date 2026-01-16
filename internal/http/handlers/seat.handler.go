package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/senatroxx/filmix-backend/internal/http/dto"
	"github.com/senatroxx/filmix-backend/internal/services"
	"github.com/senatroxx/filmix-backend/internal/utilities"
)

type SeatHandler struct {
	seatService services.ISeatService
}

func NewSeatHandler(seatService services.ISeatService) *SeatHandler {
	return &SeatHandler{seatService: seatService}
}

func (h *SeatHandler) GetSeatsForShowtime(c *fiber.Ctx) error {
	showtimeIDParam := c.Params("showtimeId")
	showtimeID, err := uuid.Parse(showtimeIDParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid showtime ID")
	}

	seats, err := h.seatService.GetSeatsForShowtime(c.Context(), showtimeID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Showtime not found")
	}

	var response []dto.SeatResponse
	for _, seat := range seats {
		resp := dto.SeatResponse{
			ID:       seat.ID,
			Row:      seat.Row,
			Number:   seat.Number,
			IsBooked: seat.IsBooked,
		}
		if seat.SeatType != nil {
			resp.SeatType = dto.SeatTypeResponse{
				ID:   seat.SeatType.ID,
				Name: seat.SeatType.Name,
			}
		}
		response = append(response, resp)
	}

	return utilities.NewSuccessResponse(c, http.StatusOK, "Seats retrieved successfully", response)
}
