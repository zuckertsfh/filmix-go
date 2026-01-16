package repositories

import (
	"database/sql"
)

type IShowtimeRepository interface {
	// Add domain methods as needed
}

type ShowtimeRepository struct {
	db *sql.DB
}

func NewShowtimeRepository(db *sql.DB) IShowtimeRepository {
	return &ShowtimeRepository{db: db}
}
