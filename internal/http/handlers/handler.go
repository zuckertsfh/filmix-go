package handlers

import "github.com/senatroxx/filmix-backend/internal/services"

type Handlers struct {
	Auth     *AuthHandler
	Movie    *MovieHandler
	Showtime *ShowtimeHandler
	Seat     *SeatHandler
	Booking  *BookingHandler
}

func RegisterHandlers(s *services.Services) *Handlers {
	return &Handlers{
		Auth:     NewAuthHandler(s.AuthService),
		Movie:    NewMovieHandler(s.MovieService),
		Showtime: NewShowtimeHandler(s.ShowtimeService),
		Seat:     NewSeatHandler(s.SeatService),
		Booking:  NewBookingHandler(s.BookingService),
	}
}
