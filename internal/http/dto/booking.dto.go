package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateBookingRequest struct {
	ShowtimeID      uuid.UUID   `json:"showtime_id" validate:"required"`
	SeatIDs         []uuid.UUID `json:"seat_ids" validate:"required,min=1"`
	PaymentMethodID uuid.UUID   `json:"payment_method_id" validate:"required"`
}

type BookingResponse struct {
	ID            uuid.UUID         `json:"id"`
	Status        string            `json:"status"`
	InvoiceNumber *string           `json:"invoice_number,omitempty"`
	Amount        int64             `json:"amount"`
	ExpiredAt     time.Time         `json:"expired_at"`
	PaidAt        *time.Time        `json:"paid_at,omitempty"`
	Showtime      BookingShowtime   `json:"showtime"`
	Theater       BookingTheater    `json:"theater"`
	Seats         []BookingSeatItem `json:"seats,omitempty"`
}

type BookingShowtime struct {
	ID    uuid.UUID  `json:"id"`
	Time  time.Time  `json:"time"`
	Movie MovieBrief `json:"movie"`
}

type BookingTheater struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type BookingSeatItem struct {
	ID       uuid.UUID `json:"id"`
	Row      string    `json:"row"`
	Number   int       `json:"number"`
	SeatType string    `json:"seat_type"`
	Price    int64     `json:"price"`
}

type BookingListResponse struct {
	ID        uuid.UUID       `json:"id"`
	Status    string          `json:"status"`
	Amount    int64           `json:"amount"`
	ExpiredAt time.Time       `json:"expired_at"`
	PaidAt    *time.Time      `json:"paid_at,omitempty"`
	Showtime  BookingShowtime `json:"showtime"`
	Theater   BookingTheater  `json:"theater"`
}
