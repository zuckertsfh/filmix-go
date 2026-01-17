package handlers

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/senatroxx/filmix-backend/internal/database/entities"
	"github.com/senatroxx/filmix-backend/internal/http/dto"
	"github.com/senatroxx/filmix-backend/internal/services"
	"github.com/senatroxx/filmix-backend/internal/utilities"
)

type BookingHandler struct {
	bookingService services.IBookingService
}

func NewBookingHandler(bookingService services.IBookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService}
}

func (h *BookingHandler) CreateBooking(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user")
	}

	var req dto.CreateBookingRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if len(req.SeatIDs) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "At least one seat is required")
	}

	input := services.CreateBookingInput{
		UserID:          userID,
		ShowtimeID:      req.ShowtimeID,
		SeatIDs:         req.SeatIDs,
		PaymentMethodID: req.PaymentMethodID,
	}

	booking, err := h.bookingService.CreateBooking(c.Context(), input)
	if err != nil {
		if errors.Is(err, services.ErrSeatsNotAvailable) {
			return fiber.NewError(fiber.StatusConflict, "One or more seats are already booked")
		}
		if errors.Is(err, services.ErrShowtimeNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Showtime not found")
		}
		// Log actual error for debugging
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := h.mapBookingToResponse(booking)
	return utilities.NewSuccessResponse(c, http.StatusCreated, "Booking created successfully", response)
}

func (h *BookingHandler) GetBooking(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user")
	}

	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid booking ID")
	}

	booking, err := h.bookingService.GetBookingByID(c.Context(), id, userID)
	if err != nil {
		if errors.Is(err, services.ErrBookingNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Booking not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get booking")
	}

	response := h.mapBookingToResponse(booking)
	return utilities.NewSuccessResponse(c, http.StatusOK, "Booking retrieved successfully", response)
}

func (h *BookingHandler) GetUserBookings(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user")
	}

	bookings, err := h.bookingService.GetUserBookings(c.Context(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get bookings")
	}

	var response []dto.BookingListResponse
	for _, b := range bookings {
		resp := dto.BookingListResponse{
			ID:        b.ID,
			Status:    b.Status,
			Amount:    b.Amount,
			ExpiredAt: b.ExpiredAt,
			PaidAt:    b.PaidAt,
		}
		if b.Showtime != nil {
			resp.Showtime = dto.BookingShowtime{
				ID:   b.Showtime.ID,
				Time: b.Showtime.Time,
			}
			if b.Showtime.Movie != nil {
				resp.Showtime.Movie = dto.MovieBrief{
					ID:        b.Showtime.Movie.ID,
					Title:     b.Showtime.Movie.Title,
					PosterURL: b.Showtime.Movie.PosterURL,
				}
			}
		}
		if b.Theater != nil {
			resp.Theater = dto.BookingTheater{
				ID:   b.Theater.ID,
				Name: b.Theater.Name,
			}
		}
		response = append(response, resp)
	}

	return utilities.NewSuccessResponse(c, http.StatusOK, "Bookings retrieved successfully", response)
}

func (h *BookingHandler) getUserID(c *fiber.Ctx) (uuid.UUID, error) {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil, errors.New("invalid user_id in token")
	}
	return uuid.Parse(userIDStr)
}

func (h *BookingHandler) mapBookingToResponse(b *entities.Transaction) dto.BookingResponse {
	resp := dto.BookingResponse{
		ID:            b.ID,
		Status:        b.Status,
		InvoiceNumber: b.InvoiceNumber,
		Amount:        b.Amount,
		ExpiredAt:     b.ExpiredAt,
		PaidAt:        b.PaidAt,
	}

	if b.Showtime != nil {
		resp.Showtime = dto.BookingShowtime{
			ID:   b.Showtime.ID,
			Time: b.Showtime.Time,
		}
		if b.Showtime.Movie != nil {
			resp.Showtime.Movie = dto.MovieBrief{
				ID:        b.Showtime.Movie.ID,
				Title:     b.Showtime.Movie.Title,
				PosterURL: b.Showtime.Movie.PosterURL,
			}
		}
	}

	if b.Theater != nil {
		resp.Theater = dto.BookingTheater{
			ID:   b.Theater.ID,
			Name: b.Theater.Name,
		}
	}

	for _, item := range b.Items {
		seatItem := dto.BookingSeatItem{
			ID:    item.SeatID,
			Price: item.Price,
		}
		if item.Seat != nil {
			seatItem.Row = item.Seat.Row
			seatItem.Number = item.Seat.Number
			if item.Seat.SeatType != nil {
				seatItem.SeatType = item.Seat.SeatType.Name
			}
		}
		resp.Seats = append(resp.Seats, seatItem)
	}

	return resp
}
