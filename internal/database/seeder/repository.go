package seeder

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/senatroxx/filmix-backend/internal/database/entities"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// --- Cinema ---

func (r *Repository) GetCinemaByName(ctx context.Context, name string) (*entities.Cinema, error) {
	c := &entities.Cinema{}
	query := `SELECT id, name, logo_url FROM cinemas WHERE name = $1`
	err := r.db.QueryRowContext(ctx, query, name).Scan(&c.ID, &c.Name, &c.LogoURL)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *Repository) CreateCinema(ctx context.Context, cinema *entities.Cinema) error {
	query := `INSERT INTO cinemas (id, name, logo_url) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, cinema.ID, cinema.Name, cinema.LogoURL)
	return err
}

func (r *Repository) GetTheaterByName(ctx context.Context, cinemaID uuid.UUID, name string) (*entities.Theater, error) {
	t := &entities.Theater{}
	query := `SELECT id, cinema_id, name FROM theaters WHERE cinema_id = $1 AND name = $2`
	err := r.db.QueryRowContext(ctx, query, cinemaID, name).Scan(&t.ID, &t.CinemaID, &t.Name)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *Repository) CreateTheater(ctx context.Context, theater *entities.Theater) error {
	query := `INSERT INTO theaters (id, cinema_id, name, address, latitude, longitude) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, theater.ID, theater.CinemaID, theater.Name, theater.Address, theater.Latitude, theater.Longitude)
	return err
}

func (r *Repository) GetStudioByName(ctx context.Context, theaterID uuid.UUID, name string) (*entities.Studio, error) {
	s := &entities.Studio{}
	query := `SELECT id, theater_id, name FROM studios WHERE theater_id = $1 AND name = $2`
	err := r.db.QueryRowContext(ctx, query, theaterID, name).Scan(&s.ID, &s.TheaterID, &s.Name)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *Repository) CreateStudio(ctx context.Context, studio *entities.Studio) error {
	query := `INSERT INTO studios (id, theater_id, name) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, studio.ID, studio.TheaterID, studio.Name)
	return err
}

func (r *Repository) CreateSeat(ctx context.Context, seat *entities.Seat) error {
	query := `INSERT INTO seats (id, studio_id, row, number, seat_type_id, active) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, seat.ID, seat.StudioID, seat.Row, seat.Number, seat.SeatTypeID, seat.Active)
	return err
}

func (r *Repository) GetSeatTypeByName(ctx context.Context, name string) (*entities.SeatType, error) {
	st := &entities.SeatType{}
	query := `SELECT id, name, cinema_id FROM seat_type WHERE name = $1`
	err := r.db.QueryRowContext(ctx, query, name).Scan(&st.ID, &st.Name, &st.CinemaID)
	if err != nil {
		return nil, err
	}
	return st, nil
}

func (r *Repository) CreateSeatType(ctx context.Context, st *entities.SeatType) error {
	query := `INSERT INTO seat_type (id, name, cinema_id) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, st.ID, st.Name, st.CinemaID)
	return err
}

func (r *Repository) CreateSeatPricing(ctx context.Context, sp *entities.SeatPricing) error {
	query := `INSERT INTO seat_pricings (id, price, day_type, seat_type_id, theater_id) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, sp.ID, sp.Price, sp.DayType, sp.SeatTypeID, sp.TheaterID)
	return err
}

// --- Movie ---

func (r *Repository) GetStatusByName(ctx context.Context, name string) (*entities.MovieStatus, error) {
	status := &entities.MovieStatus{}
	query := `SELECT id, status FROM movie_statuses WHERE status = $1`
	err := r.db.QueryRowContext(ctx, query, name).Scan(&status.ID, &status.Status)
	if err != nil {
		return nil, err
	}
	return status, nil
}

func (r *Repository) CreateStatus(ctx context.Context, status *entities.MovieStatus) error {
	query := `INSERT INTO movie_statuses (id, status) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, status.ID, status.Status)
	return err
}

func (r *Repository) GetRatingByName(ctx context.Context, name string) (*entities.MovieRating, error) {
	rating := &entities.MovieRating{}
	query := `SELECT id, rating FROM movie_ratings WHERE rating = $1`
	err := r.db.QueryRowContext(ctx, query, name).Scan(&rating.ID, &rating.Rating)
	if err != nil {
		return nil, err
	}
	return rating, nil
}

func (r *Repository) CreateRating(ctx context.Context, rating *entities.MovieRating) error {
	query := `INSERT INTO movie_ratings (id, rating) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, rating.ID, rating.Rating)
	return err
}

func (r *Repository) UpsertGenre(ctx context.Context, genre *entities.MovieGenre) error {
	query := `
		INSERT INTO movie_genres (id, genre)
		VALUES ($1, $2)
		ON CONFLICT (id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, genre.ID, genre.Genre)
	return err
}

func (r *Repository) CreateMovie(ctx context.Context, movie *entities.Movie) error {
	query := `
		INSERT INTO movies (id, title, tagline, overview, poster_url, backdrop_url, trailer_url, duration, popularity, movie_status_id, movie_rating_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.ExecContext(ctx, query,
		movie.ID, movie.Title, movie.Tagline, movie.Overview, movie.PosterURL, movie.BackdropURL, movie.TrailerURL, movie.Duration, movie.Popularity, movie.MovieStatusID, movie.MovieRatingID,
	)
	return err
}

// --- Showtime ---

func (r *Repository) CreateShowtime(ctx context.Context, showtime *entities.Showtime) error {
	query := `
		INSERT INTO showtimes (id, movie_id, studio_id, theater_id, time, expired_at, seat_pricing_id, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		showtime.ID, showtime.MovieID, showtime.StudioID, showtime.TheaterID, showtime.Time, showtime.ExpiredAt, showtime.SeatPricingID, showtime.Status,
	)
	return err
}
