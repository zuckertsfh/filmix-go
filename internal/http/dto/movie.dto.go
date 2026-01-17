package dto

import "github.com/google/uuid"

type MovieResponse struct {
	ID          uuid.UUID       `json:"id"`
	Title       string          `json:"title"`
	Tagline     string          `json:"tagline"`
	Overview    string          `json:"overview"`
	PosterURL   string          `json:"poster_url"`
	BackdropURL string          `json:"backdrop_url"`
	TrailerURL  string          `json:"trailer_url"`
	Duration    int             `json:"duration"`
	Popularity  int             `json:"popularity"`
	Status      string          `json:"status"`
	Rating      string          `json:"rating"`
	Genres      []GenreResponse `json:"genres,omitempty"`
}

type GenreResponse struct {
	ID    uuid.UUID `json:"id"`
	Genre string    `json:"genre"`
}
