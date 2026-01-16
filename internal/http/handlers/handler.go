package handlers

import "github.com/senatroxx/filmix-backend/internal/services"

type Handlers struct {
	Auth *AuthHandler
}

func RegisterHandlers(s *services.Services) *Handlers {
	return &Handlers{
		Auth: NewAuthHandler(s.AuthService),
	}
}
