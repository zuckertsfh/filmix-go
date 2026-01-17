package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/senatroxx/filmix-backend/internal/database/entities"
)

type IMovieRepository interface {
	FindAll(ctx context.Context, page, limit int) ([]entities.Movie, int, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Movie, error)
	FindNowPlaying(ctx context.Context, page, limit int) ([]entities.Movie, int, error)
}

type MovieRepository struct {
	db *sql.DB
}

func NewMovieRepository(db *sql.DB) IMovieRepository {
	return &MovieRepository{db: db}
}

func (r *MovieRepository) FindAll(ctx context.Context, page, limit int) ([]entities.Movie, int, error) {
	// Count total
	var total int
	countQuery := `SELECT COUNT(*) FROM movies`
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	query := `
		SELECT m.id, m.title, m.tagline, m.overview, m.poster_url, m.backdrop_url, 
		       m.trailer_url, m.duration, m.popularity, m.movie_status_id, m.movie_rating_id,
		       ms.id, ms.status, mr.id, mr.rating
		FROM movies m
		JOIN movie_statuses ms ON m.movie_status_id = ms.id
		JOIN movie_ratings mr ON m.movie_rating_id = mr.id
		ORDER BY m.popularity DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var movies []entities.Movie
	for rows.Next() {
		var movie entities.Movie
		var status entities.MovieStatus
		var rating entities.MovieRating

		err := rows.Scan(
			&movie.ID, &movie.Title, &movie.Tagline, &movie.Overview,
			&movie.PosterURL, &movie.BackdropURL, &movie.TrailerURL,
			&movie.Duration, &movie.Popularity, &movie.MovieStatusID, &movie.MovieRatingID,
			&status.ID, &status.Status, &rating.ID, &rating.Rating,
		)
		if err != nil {
			return nil, 0, err
		}

		movie.Status = &status
		movie.Rating = &rating
		movies = append(movies, movie)
	}

	return movies, total, nil
}

func (r *MovieRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Movie, error) {
	query := `
		SELECT m.id, m.title, m.tagline, m.overview, m.poster_url, m.backdrop_url, 
		       m.trailer_url, m.duration, m.popularity, m.movie_status_id, m.movie_rating_id,
		       ms.id, ms.status, mr.id, mr.rating
		FROM movies m
		JOIN movie_statuses ms ON m.movie_status_id = ms.id
		JOIN movie_ratings mr ON m.movie_rating_id = mr.id
		WHERE m.id = $1
	`

	var movie entities.Movie
	var status entities.MovieStatus
	var rating entities.MovieRating

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&movie.ID, &movie.Title, &movie.Tagline, &movie.Overview,
		&movie.PosterURL, &movie.BackdropURL, &movie.TrailerURL,
		&movie.Duration, &movie.Popularity, &movie.MovieStatusID, &movie.MovieRatingID,
		&status.ID, &status.Status, &rating.ID, &rating.Rating,
	)
	if err != nil {
		return nil, err
	}

	movie.Status = &status
	movie.Rating = &rating

	// Fetch genres
	genreQuery := `
		SELECT mg.id, mg.genre
		FROM movie_genres mg
		JOIN genre_movie gm ON mg.id = gm.movie_genre_id
		WHERE gm.movie_id = $1
	`
	genreRows, err := r.db.QueryContext(ctx, genreQuery, id)
	if err != nil {
		return nil, err
	}
	defer genreRows.Close()

	for genreRows.Next() {
		var genre entities.MovieGenre
		if err := genreRows.Scan(&genre.ID, &genre.Genre); err != nil {
			return nil, err
		}
		movie.Genres = append(movie.Genres, genre)
	}

	return &movie, nil
}

func (r *MovieRepository) FindNowPlaying(ctx context.Context, page, limit int) ([]entities.Movie, int, error) {
	// Count total now playing
	var total int
	countQuery := `
		SELECT COUNT(*) FROM movies m
		JOIN movie_statuses ms ON m.movie_status_id = ms.id
		WHERE ms.status = 'Now Showing'
	`
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	query := `
		SELECT m.id, m.title, m.tagline, m.overview, m.poster_url, m.backdrop_url, 
		       m.trailer_url, m.duration, m.popularity, m.movie_status_id, m.movie_rating_id,
		       ms.id, ms.status, mr.id, mr.rating
		FROM movies m
		JOIN movie_statuses ms ON m.movie_status_id = ms.id
		JOIN movie_ratings mr ON m.movie_rating_id = mr.id
		WHERE ms.status = 'Now Showing'
		ORDER BY m.popularity DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var movies []entities.Movie
	for rows.Next() {
		var movie entities.Movie
		var status entities.MovieStatus
		var rating entities.MovieRating

		err := rows.Scan(
			&movie.ID, &movie.Title, &movie.Tagline, &movie.Overview,
			&movie.PosterURL, &movie.BackdropURL, &movie.TrailerURL,
			&movie.Duration, &movie.Popularity, &movie.MovieStatusID, &movie.MovieRatingID,
			&status.ID, &status.Status, &rating.ID, &rating.Rating,
		)
		if err != nil {
			return nil, 0, err
		}

		movie.Status = &status
		movie.Rating = &rating
		movies = append(movies, movie)
	}

	return movies, total, nil
}
