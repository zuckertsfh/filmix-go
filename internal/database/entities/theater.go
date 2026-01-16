package entities

import "github.com/google/uuid"

type Theater struct {
    ID       uuid.UUID `json:"id"`
    Name     string    `json:"name"`
    Address  string    `json:"address"`
    Latitude  float64   `json:"latitude"`
    Longitude float64   `json:"longitude"`
    CinemaID uuid.UUID `json:"cinema_id"`

    Cinema   *Cinema   `json:"cinema,omitempty"`
    Studios  []Studio  `json:"studios,omitempty"`
    Pricings []SeatPricing `json:"pricings,omitempty"`
}
