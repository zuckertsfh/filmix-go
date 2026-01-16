package entities

import (
    "github.com/google/uuid"
)

type Movie struct {
    ID            uuid.UUID `json:"id"`
    Title         string    `json:"title"`
    Tagline       string    `json:"tagline"`
    Overview      string    `json:"overview"`
    PosterURL     string    `json:"poster_url"`
    BackdropURL   string    `json:"backdrop_url"`
    TrailerURL    string    `json:"trailer_url"`
    Duration      int       `json:"duration"`
    Popularity    int       `json:"popularity"`
    MovieStatusID uuid.UUID `json:"movie_status_id"`
    MovieRatingID uuid.UUID `json:"movie_rating_id"`

    Status  *MovieStatus   `json:"status,omitempty"`
    Rating  *MovieRating   `json:"rating,omitempty"`
    Genres  []MovieGenre   `json:"genres,omitempty"`
    Shows   []Showtime     `json:"showtimes,omitempty"`
}
