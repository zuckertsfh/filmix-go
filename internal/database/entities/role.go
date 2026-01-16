package entities

import "github.com/google/uuid"

type Role struct {
    ID   uuid.UUID `json:"id"`
    Name string    `json:"name"`

    Users []User `json:"users,omitempty"`
}
