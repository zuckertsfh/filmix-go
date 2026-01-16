package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/senatroxx/filmix-backend/internal/database/entities"
	"github.com/senatroxx/filmix-backend/internal/http/dto"
	"github.com/senatroxx/filmix-backend/internal/repositories"
	"github.com/senatroxx/filmix-backend/internal/utilities"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error)
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

func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		// Log error if needed, but return generic error for security
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	tokenPair, err := utilities.GenerateTokenPair(user.ID, "user")
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	tokenMetadata, err := utilities.ExtractTokenMetadata(req.RefreshToken, true)
	if err != nil {
		return nil, utilities.ErrInvalidToken
	}

	userID, err := uuid.Parse(tokenMetadata.UserID)
	if err != nil {
		return nil, utilities.ErrInvalidToken
	}

	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return nil, utilities.ErrInvalidToken
	}

	tokenPair, err := utilities.GenerateTokenPair(user.ID, "user")
	if err != nil {
		return nil, err
	}

	return &dto.RefreshTokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}
