package handlers

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if errMsg := utilities.ValidateStruct(req); errMsg != "" {
		return fiber.NewError(fiber.StatusBadRequest, errMsg)
	}

	res, err := h.authService.Login(c.Context(), req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Login failed")
	}

	return utilities.NewSuccessResponse(c, http.StatusOK, "Login successful", res)
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	req := new(dto.RegisterRequest)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if errMsg := utilities.ValidateStruct(req); errMsg != "" {
		return fiber.NewError(fiber.StatusBadRequest, errMsg)
	}

	res, err := h.authService.Register(c.Context(), req)
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyRegistered) {
			return fiber.NewError(fiber.StatusConflict, err.Error())
		}
		if errors.Is(err, services.ErrRoleNotFound) {
			return fiber.NewError(fiber.StatusInternalServerError, "Role configuration error")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return utilities.NewSuccessResponse(c, http.StatusCreated, "User registered successfully", res)
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Missing Authorization header")
	}

	if len(authHeader) <= 7 || authHeader[:7] != "Bearer " {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid Authorization header format")
	}

	refreshToken := authHeader[7:]
	req := &dto.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	res, err := h.authService.RefreshToken(c.Context(), req)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}

	return utilities.NewSuccessResponse(c, http.StatusOK, "Token refreshed successfully", res)
}
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	// gofiber/contrib/jwt stores claims in c.Locals("user") by default
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user ID in token")
	}

	res, err := h.authService.GetProfile(c.Context(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	return utilities.NewSuccessResponse(c, http.StatusOK, "User profile retrieved successfully", res)
}
