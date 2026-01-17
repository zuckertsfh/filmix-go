package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/senatroxx/filmix-backend/internal/database/entities"
	"github.com/senatroxx/filmix-backend/internal/repositories"
)

type SeatWithAvailability struct {
	entities.Seat
	IsBooked bool `json:"is_booked"`
}

type ISeatService interface {
	GetSeatsForShowtime(ctx context.Context, showtimeID uuid.UUID) ([]SeatWithAvailability, error)
}

type SeatService struct {
	seatRepo     repositories.ISeatRepository
	showtimeRepo repositories.IShowtimeRepository
}

func NewSeatService(seatRepo repositories.ISeatRepository, showtimeRepo repositories.IShowtimeRepository) ISeatService {
	return &SeatService{
		seatRepo:     seatRepo,
		showtimeRepo: showtimeRepo,
	}
}

func (s *SeatService) GetSeatsForShowtime(ctx context.Context, showtimeID uuid.UUID) ([]SeatWithAvailability, error) {
	showtime, err := s.showtimeRepo.FindByID(ctx, showtimeID)
	if err != nil {
		return nil, err
	}

	seats, err := s.seatRepo.FindByStudioID(ctx, showtime.StudioID)
	if err != nil {
		return nil, err
	}

	bookedIDs, err := s.seatRepo.FindBookedSeatIDs(ctx, showtimeID)
	if err != nil {
		return nil, err
	}

	bookedMap := make(map[uuid.UUID]bool)
	for _, id := range bookedIDs {
		bookedMap[id] = true
	}

	var result []SeatWithAvailability
	for _, seat := range seats {
		result = append(result, SeatWithAvailability{
			Seat:     seat,
			IsBooked: bookedMap[seat.ID],
		})
	}

	return result, nil
}
