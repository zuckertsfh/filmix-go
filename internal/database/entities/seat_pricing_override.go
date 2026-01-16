package entities

import "github.com/google/uuid"

type SeatPricingOverride struct {
	ID         uuid.UUID `json:"id"`
	Price      int64     `json:"price"`
	Notes      string    `json:"notes"`
	MovieID    uuid.UUID `json:"movie_id"`
	SeatTypeID uuid.UUID `json:"seat_type_id"`
	TheaterID  uuid.UUID `json:"theater_id"`

	Movie    *Movie    `json:"movie,omitempty"`
	SeatType *SeatType `json:"seat_type,omitempty"`
	Theater  *Theater  `json:"theater,omitempty"`
}
