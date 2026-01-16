package entities

import "github.com/google/uuid"

type Cinema struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	LogoURL string    `json:"logo_url"`

	Theaters  []Theater  `json:"theaters,omitempty"`
	SeatTypes []SeatType `json:"seat_types,omitempty"`
}
