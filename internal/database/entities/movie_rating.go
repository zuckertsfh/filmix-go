package entities

import "github.com/google/uuid"

type MovieRating struct {
    ID     uuid.UUID `json:"id"`
    Rating string    `json:"rating"`

    Movies []Movie `json:"movies,omitempty"`
}
