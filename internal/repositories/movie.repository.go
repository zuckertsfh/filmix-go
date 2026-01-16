package repositories

import (
	"database/sql"
)

type IMovieRepository interface {
	// Add domain methods as needed
}

type MovieRepository struct {
	db *sql.DB
}

func NewMovieRepository(db *sql.DB) IMovieRepository {
	return &MovieRepository{db: db}
}
