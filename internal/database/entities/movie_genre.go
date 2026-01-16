package entities

import "github.com/google/uuid"

type MovieGenre struct {
    ID    uuid.UUID `json:"id"`
    Genre string    `json:"genre"`

    Movies []Movie `json:"movies,omitempty"`
}
