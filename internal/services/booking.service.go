package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/senatroxx/filmix-backend/internal/database/entities"
	"github.com/senatroxx/filmix-backend/internal/repositories"
)

var (
	ErrSeatsNotAvailable = errors.New("one or more seats are not available")
	ErrShowtimeNotFound  = errors.New("showtime not found")
	ErrBookingNotFound   = errors.New("booking not found")
)

type CreateBookingInput struct {
	UserID          uuid.UUID
	ShowtimeID      uuid.UUID
	SeatIDs         []uuid.UUID
	PaymentMethodID uuid.UUID
}

type IBookingService interface {
	CreateBooking(ctx context.Context, input CreateBookingInput) (*entities.Transaction, error)
	GetBookingByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entities.Transaction, error)
	GetUserBookings(ctx context.Context, userID uuid.UUID) ([]entities.Transaction, error)
}

type BookingService struct {
	bookingRepo  repositories.IBookingRepository
	showtimeRepo repositories.IShowtimeRepository
	seatRepo     repositories.ISeatRepository
}

func NewBookingService(
	bookingRepo repositories.IBookingRepository,
	showtimeRepo repositories.IShowtimeRepository,
	seatRepo repositories.ISeatRepository,
) IBookingService {
	return &BookingService{
		bookingRepo:  bookingRepo,
		showtimeRepo: showtimeRepo,
		seatRepo:     seatRepo,
	}
}

func (s *BookingService) CreateBooking(ctx context.Context, input CreateBookingInput) (*entities.Transaction, error) {
	showtime, err := s.showtimeRepo.FindByID(ctx, input.ShowtimeID)
	if err != nil {
		return nil, ErrShowtimeNotFound
	}

	available, err := s.bookingRepo.CheckSeatsAvailable(ctx, input.ShowtimeID, input.SeatIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to check seat availability: %w", err)
	}
	if !available {
		return nil, ErrSeatsNotAvailable
	}

	seats, err := s.seatRepo.FindByStudioID(ctx, showtime.StudioID)
	if err != nil {
		return nil, fmt.Errorf("failed to get seats: %w", err)
	}

	seatMap := make(map[uuid.UUID]entities.Seat)
	for _, seat := range seats {
		seatMap[seat.ID] = seat
	}

	var totalAmount int64
	var items []entities.TransactionItem
	txID := uuid.New()

	for _, seatID := range input.SeatIDs {
		seat, exists := seatMap[seatID]
		if !exists {
			return nil, fmt.Errorf("seat %s not found in studio", seatID)
		}

		price := showtime.Pricing.Price

		items = append(items, entities.TransactionItem{
			ID:            uuid.New(),
			Price:         price,
			TransactionID: txID,
			SeatID:        seatID,
			SeatTypeID:    seat.SeatTypeID,
		})

		totalAmount += price
	}

	invoiceNumber := fmt.Sprintf("INV-%s", txID.String()[:8])
	tx := &entities.Transaction{
		ID:              txID,
		Status:          "pending",
		ExternalRef:     "",
		InvoiceNumber:   &invoiceNumber,
		Amount:          totalAmount,
		ExpiredAt:       time.Now().Add(15 * time.Minute), // 15 minutes to pay
		PaymentMethodID: input.PaymentMethodID,
		ShowtimeID:      input.ShowtimeID,
		TheaterID:       showtime.TheaterID,
		UserID:          input.UserID,
	}

	err = s.bookingRepo.Create(ctx, tx, items)
	if err != nil {
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}

	return s.bookingRepo.FindByID(ctx, txID)
}

func (s *BookingService) GetBookingByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entities.Transaction, error) {
	booking, err := s.bookingRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrBookingNotFound
	}

	if booking.UserID != userID {
		return nil, ErrBookingNotFound
	}

	return booking, nil
}

func (s *BookingService) GetUserBookings(ctx context.Context, userID uuid.UUID) ([]entities.Transaction, error) {
	return s.bookingRepo.FindByUserID(ctx, userID)
}
