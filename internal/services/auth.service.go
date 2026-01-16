package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/senatroxx/filmix-backend/internal/database/entities"
	"github.com/senatroxx/filmix-backend/internal/http/dto"
	"github.com/senatroxx/filmix-backend/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
}

type AuthService struct {
	userRepository repositories.IUserRepository
}

func NewAuthService(userRepo repositories.IUserRepository) IAuthService {
	return &AuthService{userRepository: userRepo}
}

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	existingUser, _ := s.userRepository.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, ErrEmailAlreadyRegistered
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	userRole, err := s.userRepository.GetRoleByName(ctx, "user")
	if err != nil {
		return nil, ErrRoleNotFound
	}

	newUser := &entities.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		RoleID:   userRole.ID,
	}

	err = s.userRepository.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return &dto.RegisterResponse{
		ID:    newUser.ID,
		Name:  newUser.Name,
		Email: newUser.Email,
	}, nil
}
