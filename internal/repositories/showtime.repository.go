package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/senatroxx/filmix-backend/internal/database/entities"
)

type IShowtimeRepository interface {
	FindByMovieID(ctx context.Context, movieID uuid.UUID, date *time.Time) ([]entities.Showtime, error)
	FindByTheaterID(ctx context.Context, theaterID uuid.UUID, date *time.Time) ([]entities.Showtime, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Showtime, error)
}

type ShowtimeRepository struct {
	db *sql.DB
}

func NewShowtimeRepository(db *sql.DB) IShowtimeRepository {
	return &ShowtimeRepository{db: db}
}

func (r *ShowtimeRepository) FindByMovieID(ctx context.Context, movieID uuid.UUID, date *time.Time) ([]entities.Showtime, error) {
	query := `
		SELECT 
			s.id, s.status, s.time, s.expired_at, s.movie_id, s.studio_id, s.theater_id, s.seat_pricing_id,
			st.id, st.name, st.theater_id,
			t.id, t.name, t.address, t.latitude, t.longitude, t.cinema_id,
			c.id, c.name, c.logo_url,
			sp.id, sp.price, sp.day_type
		FROM showtimes s
		JOIN studios st ON s.studio_id = st.id
		JOIN theaters t ON s.theater_id = t.id
		JOIN cinemas c ON t.cinema_id = c.id
		JOIN seat_pricings sp ON s.seat_pricing_id = sp.id
		WHERE s.movie_id = $1 AND s.status = true AND s.time > NOW()
	`

	args := []interface{}{movieID}
	if date != nil {
		query += ` AND DATE(s.time) = $2`
		args = append(args, date.Format("2006-01-02"))
	}

	query += ` ORDER BY s.time ASC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanShowtimes(rows)
}

func (r *ShowtimeRepository) FindByTheaterID(ctx context.Context, theaterID uuid.UUID, date *time.Time) ([]entities.Showtime, error) {
	query := `
		SELECT 
			s.id, s.status, s.time, s.expired_at, s.movie_id, s.studio_id, s.theater_id, s.seat_pricing_id,
			st.id, st.name, st.theater_id,
			t.id, t.name, t.address, t.latitude, t.longitude, t.cinema_id,
			c.id, c.name, c.logo_url,
			sp.id, sp.price, sp.day_type,
			m.id, m.title, m.poster_url, m.duration
		FROM showtimes s
		JOIN studios st ON s.studio_id = st.id
		JOIN theaters t ON s.theater_id = t.id
		JOIN cinemas c ON t.cinema_id = c.id
		JOIN seat_pricings sp ON s.seat_pricing_id = sp.id
		JOIN movies m ON s.movie_id = m.id
		WHERE s.theater_id = $1 AND s.status = true AND s.time > NOW()
	`

	args := []interface{}{theaterID}
	if date != nil {
		query += ` AND DATE(s.time) = $2`
		args = append(args, date.Format("2006-01-02"))
	}

	query += ` ORDER BY s.time ASC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanShowtimesWithMovie(rows)
}

func (r *ShowtimeRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Showtime, error) {
	query := `
		SELECT 
			s.id, s.status, s.time, s.expired_at, s.movie_id, s.studio_id, s.theater_id, s.seat_pricing_id,
			st.id, st.name, st.theater_id,
			t.id, t.name, t.address, t.latitude, t.longitude, t.cinema_id,
			c.id, c.name, c.logo_url,
			sp.id, sp.price, sp.day_type,
			m.id, m.title, m.poster_url, m.duration
		FROM showtimes s
		JOIN studios st ON s.studio_id = st.id
		JOIN theaters t ON s.theater_id = t.id
		JOIN cinemas c ON t.cinema_id = c.id
		JOIN seat_pricings sp ON s.seat_pricing_id = sp.id
		JOIN movies m ON s.movie_id = m.id
		WHERE s.id = $1
	`

	var showtime entities.Showtime
	var studio entities.Studio
	var theater entities.Theater
	var cinema entities.Cinema
	var pricing entities.SeatPricing
	var movie entities.Movie

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&showtime.ID, &showtime.Status, &showtime.Time, &showtime.ExpiredAt,
		&showtime.MovieID, &showtime.StudioID, &showtime.TheaterID, &showtime.SeatPricingID,
		&studio.ID, &studio.Name, &studio.TheaterID,
		&theater.ID, &theater.Name, &theater.Address, &theater.Latitude, &theater.Longitude, &theater.CinemaID,
		&cinema.ID, &cinema.Name, &cinema.LogoURL,
		&pricing.ID, &pricing.Price, &pricing.DayType,
		&movie.ID, &movie.Title, &movie.PosterURL, &movie.Duration,
	)
	if err != nil {
		return nil, err
	}

	theater.Cinema = &cinema
	showtime.Studio = &studio
	showtime.Theater = &theater
	showtime.Pricing = &pricing
	showtime.Movie = &movie

	return &showtime, nil
}

func (r *ShowtimeRepository) scanShowtimes(rows *sql.Rows) ([]entities.Showtime, error) {
	var showtimes []entities.Showtime

	for rows.Next() {
		var showtime entities.Showtime
		var studio entities.Studio
		var theater entities.Theater
		var cinema entities.Cinema
		var pricing entities.SeatPricing

		err := rows.Scan(
			&showtime.ID, &showtime.Status, &showtime.Time, &showtime.ExpiredAt,
			&showtime.MovieID, &showtime.StudioID, &showtime.TheaterID, &showtime.SeatPricingID,
			&studio.ID, &studio.Name, &studio.TheaterID,
			&theater.ID, &theater.Name, &theater.Address, &theater.Latitude, &theater.Longitude, &theater.CinemaID,
			&cinema.ID, &cinema.Name, &cinema.LogoURL,
			&pricing.ID, &pricing.Price, &pricing.DayType,
		)
		if err != nil {
			return nil, err
		}

		theater.Cinema = &cinema
		showtime.Studio = &studio
		showtime.Theater = &theater
		showtime.Pricing = &pricing

		showtimes = append(showtimes, showtime)
	}

	return showtimes, nil
}

func (r *ShowtimeRepository) scanShowtimesWithMovie(rows *sql.Rows) ([]entities.Showtime, error) {
	var showtimes []entities.Showtime

	for rows.Next() {
		var showtime entities.Showtime
		var studio entities.Studio
		var theater entities.Theater
		var cinema entities.Cinema
		var pricing entities.SeatPricing
		var movie entities.Movie

		err := rows.Scan(
			&showtime.ID, &showtime.Status, &showtime.Time, &showtime.ExpiredAt,
			&showtime.MovieID, &showtime.StudioID, &showtime.TheaterID, &showtime.SeatPricingID,
			&studio.ID, &studio.Name, &studio.TheaterID,
			&theater.ID, &theater.Name, &theater.Address, &theater.Latitude, &theater.Longitude, &theater.CinemaID,
			&cinema.ID, &cinema.Name, &cinema.LogoURL,
			&pricing.ID, &pricing.Price, &pricing.DayType,
			&movie.ID, &movie.Title, &movie.PosterURL, &movie.Duration,
		)
		if err != nil {
			return nil, err
		}

		theater.Cinema = &cinema
		showtime.Studio = &studio
		showtime.Theater = &theater
		showtime.Pricing = &pricing
		showtime.Movie = &movie

		showtimes = append(showtimes, showtime)
	}

	return showtimes, nil
}
