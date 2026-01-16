package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/senatroxx/filmix-backend/internal/http/dto"
	"github.com/senatroxx/filmix-backend/internal/services"
	"github.com/senatroxx/filmix-backend/internal/utilities"
)

type AuthHandler struct {
	authService services.IAuthService
}

func NewAuthHandler(authService services.IAuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	return utilities.NewSuccessResponse(c, http.StatusOK, "Login successful", nil)
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	req := new(dto.RegisterRequest)

	// Parse Body
	if err := c.BodyParser(req); err != nil {
		return utilities.NewErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	// TODO: Add Validator here
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return utilities.NewErrorResponse(c, http.StatusBadRequest, "Missing required fields")
	}

	// Call Service
	res, err := h.authService.Register(c.Context(), req)
	if err != nil {
		// Example error handling
		if err.Error() == "email already registered" {
			return utilities.NewErrorResponse(c, http.StatusConflict, err.Error())
		}
		if err.Error() == "default role 'user' not found" {
			return utilities.NewErrorResponse(c, http.StatusInternalServerError, "Role configuration error")
		}
		return utilities.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utilities.NewSuccessResponse(c, http.StatusCreated, "User registered successfully", res)
}
