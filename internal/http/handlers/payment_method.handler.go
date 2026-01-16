package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/senatroxx/filmix-backend/internal/utilities"
)

type PaymentMethodHandler struct {
	db *sql.DB
}

func NewPaymentMethodHandler(db *sql.DB) *PaymentMethodHandler {
	return &PaymentMethodHandler{db: db}
}

type PaymentMethodResponse struct {
	ID      uuid.UUID `json:"id"`
	Code    string    `json:"code"`
	Name    string    `json:"name"`
	LogoURL string    `json:"logo_url"`
}

func (h *PaymentMethodHandler) GetPaymentMethods(c *fiber.Ctx) error {
	query := `SELECT id, code, name, logo_url FROM payment_methods WHERE active = true`

	rows, err := h.db.QueryContext(c.Context(), query)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get payment methods")
	}
	defer rows.Close()

	var methods []PaymentMethodResponse
	for rows.Next() {
		var m PaymentMethodResponse
		if err := rows.Scan(&m.ID, &m.Code, &m.Name, &m.LogoURL); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to scan payment method")
		}
		methods = append(methods, m)
	}

	return utilities.NewSuccessResponse(c, http.StatusOK, "Payment methods retrieved successfully", methods)
}
