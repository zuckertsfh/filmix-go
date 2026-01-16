package services

import "github.com/senatroxx/filmix-backend/internal/repositories"

type Services struct {
	AuthService IAuthService
}

func RegisterServices(r *repositories.Repositories) *Services {
	return &Services{
		AuthService: NewAuthService(r.UserRepository),
	}
}
