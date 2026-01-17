package repositories

import (
	"database/sql"
)

type Repositories struct {
	UserRepository     IUserRepository
	MovieRepository    IMovieRepository
	CinemaRepository   ICinemaRepository
	ShowtimeRepository IShowtimeRepository
	SeatRepository     ISeatRepository
	BookingRepository  IBookingRepository
}

func RegisterRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		UserRepository:     NewUserRepository(db),
		MovieRepository:    NewMovieRepository(db),
		CinemaRepository:   NewCinemaRepository(db),
		ShowtimeRepository: NewShowtimeRepository(db),
		SeatRepository:     NewSeatRepository(db),
		BookingRepository:  NewBookingRepository(db),
	}
}
