package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/senatroxx/filmix-backend/internal/database/entities"
	"github.com/senatroxx/filmix-backend/internal/repositories"
)

type IShowtimeService interface {
	GetShowtimesByMovieID(ctx context.Context, movieID uuid.UUID, date *time.Time) ([]entities.Showtime, error)
	GetShowtimesByTheaterID(ctx context.Context, theaterID uuid.UUID, date *time.Time) ([]entities.Showtime, error)
	GetShowtimeByID(ctx context.Context, id uuid.UUID) (*entities.Showtime, error)
}

type ShowtimeService struct {
	showtimeRepo repositories.IShowtimeRepository
}

func NewShowtimeService(showtimeRepo repositories.IShowtimeRepository) IShowtimeService {
	return &ShowtimeService{showtimeRepo: showtimeRepo}
}

func (s *ShowtimeService) GetShowtimesByMovieID(ctx context.Context, movieID uuid.UUID, date *time.Time) ([]entities.Showtime, error) {
	return s.showtimeRepo.FindByMovieID(ctx, movieID, date)
}

func (s *ShowtimeService) GetShowtimesByTheaterID(ctx context.Context, theaterID uuid.UUID, date *time.Time) ([]entities.Showtime, error) {
	return s.showtimeRepo.FindByTheaterID(ctx, theaterID, date)
}

func (s *ShowtimeService) GetShowtimeByID(ctx context.Context, id uuid.UUID) (*entities.Showtime, error) {
	return s.showtimeRepo.FindByID(ctx, id)
}
