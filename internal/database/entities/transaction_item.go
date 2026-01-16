package entities

import "github.com/google/uuid"

type TransactionItem struct {
    ID            uuid.UUID `json:"id"`
    Price         int64     `json:"price"`
    TransactionID uuid.UUID `json:"transaction_id"`
    SeatID        uuid.UUID `json:"seat_id"`
    SeatTypeID    uuid.UUID `json:"seat_type_id"`

    Transaction *Transaction `json:"transaction,omitempty"`
    Seat        *Seat        `json:"seat,omitempty"`
    SeatType    *SeatType    `json:"seat_type,omitempty"`
}
