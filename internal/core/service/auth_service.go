package service

import (
	"context"

	"github.com/ecelayes/pms-core/internal/core/domain"
	"github.com/ecelayes/pms-core/internal/core/ports"
	"github.com/ecelayes/pms-core/pkg/auth"
)

type AuthService struct {
	repo ports.AuthRepository
}

func NewAuthService(repo ports.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(ctx context.Context, req domain.RegisterRequest) (string, error) {
	existingUser, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", err
	}
	if existingUser != nil {
		return "", domain.ErrEmailAlreadyExists
	}
	
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return "", err
	}

	return s.repo.RegisterUser(ctx, req.Email, hash)
}

func (s *AuthService) Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	
	if user == nil {
		return nil, domain.ErrInvalidCredentials
	}

	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		return nil, domain.ErrInvalidCredentials
	}

	hotelID := ""
	if user.HotelID != nil {
		hotelID = *user.HotelID
	}

	token, err := auth.GenerateToken(user.ID, hotelID, user.Role)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{Token: token}, nil
}
