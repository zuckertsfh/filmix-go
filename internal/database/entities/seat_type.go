package entities

import "github.com/google/uuid"

type SeatType struct {
    ID       uuid.UUID `json:"id"`
    Name     string    `json:"name"`
    CinemaID uuid.UUID `json:"cinema_id"`

    Cinema   *Cinema       `json:"cinema,omitempty"`
    Seats    []Seat        `json:"seats,omitempty"`
    Pricings []SeatPricing `json:"pricings,omitempty"`
}
