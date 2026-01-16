package entities

import "github.com/google/uuid"

type PaymentMethodType struct {
    ID   uuid.UUID `json:"id"`
    Name string    `json:"name"`

    Methods []PaymentMethod `json:"methods,omitempty"`
}
