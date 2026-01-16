package entities

import "github.com/google/uuid"

type Seat struct {
    ID         uuid.UUID `json:"id"`
    Row        string    `json:"row"`
    Number     int       `json:"number"`
    Active     bool      `json:"active"`
    StudioID   uuid.UUID `json:"studio_id"`
    SeatTypeID uuid.UUID `json:"seat_type_id"`

    Studio   *Studio   `json:"studio,omitempty"`
    SeatType *SeatType `json:"seat_type,omitempty"`
}
