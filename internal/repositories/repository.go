package repositories

import (
	"database/sql"
)

type Repositories struct {
	UserRepository     IUserRepository
	MovieRepository    IMovieRepository
	CinemaRepository   ICinemaRepository
	ShowtimeRepository IShowtimeRepository
}

func RegisterRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		UserRepository:     NewUserRepository(db),
		MovieRepository:    NewMovieRepository(db),
		CinemaRepository:   NewCinemaRepository(db),
		ShowtimeRepository: NewShowtimeRepository(db),
	}
}
