package dto

import "github.com/google/uuid"

type SeatResponse struct {
	ID       uuid.UUID        `json:"id"`
	Row      string           `json:"row"`
	Number   int              `json:"number"`
	SeatType SeatTypeResponse `json:"seat_type"`
	IsBooked bool             `json:"is_booked"`
}

type SeatTypeResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
