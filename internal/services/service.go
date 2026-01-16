package services

import "github.com/senatroxx/filmix-backend/internal/repositories"

type Services struct {
	AuthService     IAuthService
	MovieService    IMovieService
	ShowtimeService IShowtimeService
}

func RegisterServices(r *repositories.Repositories) *Services {
	return &Services{
		AuthService:     NewAuthService(r.UserRepository),
		MovieService:    NewMovieService(r.MovieRepository),
		ShowtimeService: NewShowtimeService(r.ShowtimeRepository),
	}
}
