package entities

import (
    "time"
    "github.com/google/uuid"
)

type Transaction struct {
    ID             uuid.UUID  `json:"id"`
    Status         string     `json:"status"`
    ExternalRef    string     `json:"external_ref"`
    InvoiceNumber  *string    `json:"invoice_number,omitempty"`
    Amount         int64      `json:"amount"`
    ExpiredAt      time.Time  `json:"expired_at"`
    PaidAt         *time.Time `json:"paid_at,omitempty"`
    PaymentMethodID uuid.UUID `json:"payment_method_id"`
    ShowtimeID      uuid.UUID `json:"showtime_id"`
    TheaterID       uuid.UUID `json:"theater_id"`
    UserID          uuid.UUID `json:"user_id"`

    PaymentMethod *PaymentMethod `json:"payment_method,omitempty"`
    Showtime      *Showtime      `json:"showtime,omitempty"`
    Theater       *Theater       `json:"theater,omitempty"`
    User          *User          `json:"user,omitempty"`
    Items         []TransactionItem `json:"items,omitempty"`
}
