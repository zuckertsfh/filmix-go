package dto

import (
	"time"

	"github.com/google/uuid"
)

type ShowtimeResponse struct {
	ID        uuid.UUID       `json:"id"`
	Time      time.Time       `json:"time"`
	ExpiredAt time.Time       `json:"expired_at"`
	Studio    StudioResponse  `json:"studio"`
	Theater   TheaterResponse `json:"theater"`
	Price     int64           `json:"price"`
	Movie     *MovieBrief     `json:"movie,omitempty"`
}

type StudioResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type TheaterResponse struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	Address   string         `json:"address"`
	Latitude  float64        `json:"latitude"`
	Longitude float64        `json:"longitude"`
	Cinema    CinemaResponse `json:"cinema"`
}

type CinemaResponse struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	LogoURL string    `json:"logo_url"`
}

type MovieBrief struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	PosterURL string    `json:"poster_url"`
	Duration  int       `json:"duration"`
}
