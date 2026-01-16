package entities

import "github.com/google/uuid"

type MovieStatus struct {
    ID     uuid.UUID `json:"id"`
    Status string    `json:"status"`

    Movies []Movie `json:"movies,omitempty"`
}
