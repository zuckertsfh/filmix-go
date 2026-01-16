package entities

import (
    "time"
    "github.com/google/uuid"
)

type Showtime struct {
    ID                    uuid.UUID  `json:"id"`
    Status                bool       `json:"status"`
    Time                  time.Time  `json:"time"`
    ExpiredAt             time.Time  `json:"expired_at"`
    MovieID               uuid.UUID  `json:"movie_id"`
    StudioID              uuid.UUID  `json:"studio_id"`
    TheaterID             uuid.UUID  `json:"theater_id"`
    SeatPricingID         uuid.UUID  `json:"seat_pricing_id"`
    SeatPricingOverrideID *uuid.UUID `json:"seat_pricing_override_id,omitempty"`

    Movie    *Movie              `json:"movie,omitempty"`
    Studio   *Studio             `json:"studio,omitempty"`
    Theater  *Theater            `json:"theater,omitempty"`
    Pricing  *SeatPricing        `json:"pricing,omitempty"`
    Override *SeatPricingOverride `json:"override,omitempty"`
}
