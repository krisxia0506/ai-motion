package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/xiajiayi/ai-motion/internal/application/dto"
	"github.com/xiajiayi/ai-motion/internal/domain/user"
	"github.com/xiajiayi/ai-motion/internal/infrastructure/auth"
)

type AuthService struct {
	userRepo   user.UserRepository
	jwtService *auth.JWTService
}

func NewAuthService(userRepo user.UserRepository, jwtService *auth.JWTService) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, user.ErrUserAlreadyExists
	}

	existingUser, err = s.userRepo.FindByUsername(ctx, req.Username)
	if err == nil && existingUser != nil {
		return nil, errors.New("username already exists")
	}

	newUser, err := user.NewUser(req.Username, req.Email, req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if err := s.userRepo.Save(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	token, err := s.jwtService.GenerateToken(string(newUser.ID), newUser.Username, newUser.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        string(newUser.ID),
			Username:  newUser.Username,
			Email:     newUser.Email,
			CreatedAt: newUser.CreatedAt,
		},
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	u, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, user.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if err := u.VerifyPassword(req.Password); err != nil {
		return nil, err
	}

	token, err := s.jwtService.GenerateToken(string(u.ID), u.Username, u.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        string(u.ID),
			Username:  u.Username,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
		},
	}, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*dto.UserResponse, error) {
	u, err := s.userRepo.FindByID(ctx, user.UserID(userID))
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        string(u.ID),
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}, nil
}
