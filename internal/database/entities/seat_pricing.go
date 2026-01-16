package entities

import "github.com/google/uuid"

type SeatPricing struct {
	ID         uuid.UUID `json:"id"`
	Price      int64     `json:"price"`
	DayType    string    `json:"day_type"`
	SeatTypeID uuid.UUID `json:"seat_type_id"`
	TheaterID  uuid.UUID `json:"theater_id"`

	SeatType *SeatType `json:"seat_type,omitempty"`
	Theater  *Theater  `json:"theater,omitempty"`
}
