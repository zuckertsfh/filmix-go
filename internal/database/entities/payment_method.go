package entities

import "github.com/google/uuid"

type PaymentMethod struct {
    ID                  uuid.UUID `json:"id"`
    Code                string    `json:"code"`
    Name                string    `json:"name"`
    LogoURL             string    `json:"logo_url"`
    Active              bool      `json:"active"`
    PaymentMethodTypeID uuid.UUID `json:"payment_method_type_id"`

    Type *PaymentMethodType `json:"type,omitempty"`
}
