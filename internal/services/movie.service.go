package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/senatroxx/filmix-backend/internal/database/entities"
	"github.com/senatroxx/filmix-backend/internal/repositories"
)

type IMovieService interface {
	GetAllMovies(ctx context.Context, page, limit int) ([]entities.Movie, int, error)
	GetMovieByID(ctx context.Context, id uuid.UUID) (*entities.Movie, error)
	GetNowPlayingMovies(ctx context.Context, page, limit int) ([]entities.Movie, int, error)
}

type MovieService struct {
	movieRepo repositories.IMovieRepository
}

func NewMovieService(movieRepo repositories.IMovieRepository) IMovieService {
	return &MovieService{movieRepo: movieRepo}
}

func (s *MovieService) GetAllMovies(ctx context.Context, page, limit int) ([]entities.Movie, int, error) {
	return s.movieRepo.FindAll(ctx, page, limit)
}

func (s *MovieService) GetMovieByID(ctx context.Context, id uuid.UUID) (*entities.Movie, error) {
	return s.movieRepo.FindByID(ctx, id)
}

func (s *MovieService) GetNowPlayingMovies(ctx context.Context, page, limit int) ([]entities.Movie, int, error) {
	return s.movieRepo.FindNowPlaying(ctx, page, limit)
}
