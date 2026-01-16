package entities

import "github.com/google/uuid"

type Studio struct {
    ID        uuid.UUID `json:"id"`
    Name      string    `json:"name"`
    TheaterID uuid.UUID `json:"theater_id"`

    Theater *Theater `json:"theater,omitempty"`
    Seats   []Seat   `json:"seats,omitempty"`
}
