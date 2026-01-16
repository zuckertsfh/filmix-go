package handlers

import (
	"errors"
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
	req := new(dto.LoginRequest)

	if err := c.BodyParser(req); err != nil {
		return utilities.NewErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if errMsg := utilities.ValidateStruct(req); errMsg != "" {
		return utilities.NewErrorResponse(c, http.StatusBadRequest, errMsg)
	}

	res, err := h.authService.Login(c.Context(), req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			return utilities.NewErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		}
		return utilities.NewErrorResponse(c, http.StatusInternalServerError, "Login failed")
	}

	return utilities.NewSuccessResponse(c, http.StatusOK, "Login successful", res)
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	req := new(dto.RegisterRequest)

	if err := c.BodyParser(req); err != nil {
		return utilities.NewErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if errMsg := utilities.ValidateStruct(req); errMsg != "" {
		return utilities.NewErrorResponse(c, http.StatusBadRequest, errMsg)
	}

	res, err := h.authService.Register(c.Context(), req)
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyRegistered) {
			return utilities.NewErrorResponse(c, http.StatusConflict, err.Error())
		}
		if errors.Is(err, services.ErrRoleNotFound) {
			return utilities.NewErrorResponse(c, http.StatusInternalServerError, "Role configuration error")
		}
		return utilities.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utilities.NewSuccessResponse(c, http.StatusCreated, "User registered successfully", res)
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return utilities.NewErrorResponse(c, http.StatusUnauthorized, "Missing Authorization header")
	}

	if len(authHeader) <= 7 || authHeader[:7] != "Bearer " {
		return utilities.NewErrorResponse(c, http.StatusUnauthorized, "Invalid Authorization header format")
	}

	refreshToken := authHeader[7:]
	req := &dto.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	res, err := h.authService.RefreshToken(c.Context(), req)
	if err != nil {
		return utilities.NewErrorResponse(c, http.StatusUnauthorized, "Invalid refresh token")
	}

	return utilities.NewSuccessResponse(c, http.StatusOK, "Token refreshed successfully", res)
}
