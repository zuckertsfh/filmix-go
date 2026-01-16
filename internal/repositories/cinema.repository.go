package repositories

import (
	"database/sql"
)

type ICinemaRepository interface {
	// Add domain methods as needed
}

type CinemaRepository struct {
	db *sql.DB
}

func NewCinemaRepository(db *sql.DB) ICinemaRepository {
	return &CinemaRepository{db: db}
}

// Keep existing domain methods if they are used by Services.
// Currently, Services seem empty or minimal.
// I will keep the struct clean as per user request.
// The interface methods (CreateCinema, GetCinemaByName etc) were only used by Seeder.
// So removing them is correct.
