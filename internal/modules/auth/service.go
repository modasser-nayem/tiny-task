package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository   Repository
	tokenManager *TokenManager
}

func NewService(repository Repository, tokenManager *TokenManager) *Service {
	return &Service{
		repository:   repository,
		tokenManager: tokenManager,
	}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	name := strings.TrimSpace(req.Name)
	email := strings.ToLower(strings.TrimSpace(req.Email))
	password := req.Password

	if name == "" {
		return nil, errors.New("name is required")
	}

	if email == "" {
		return nil, errors.New("email is required")
	}

	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	if len(password) > 72 {
		return nil, errors.New("password cannot exceed 72 bytes")
	}

	existingUser, err := s.repository.FindByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &User{
		Name:         name,
		Email:        email,
		PasswordHash: string(passwordHash),
	}

	createdUser, err := s.repository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return &RegisterResponse{
		ID:    createdUser.ID,
		Name:  createdUser.Name,
		Email: createdUser.Email,
	}, nil
}

func (s *Service) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	if req == nil {
		return nil, errors.New("request body is required")
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" || req.Password == "" {
		return nil, errors.New("email and password are required")
	}

	user, err := s.repository.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := s.tokenManager.Generate(user.ID)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	return &LoginResponse{
		AccessToken: accessToken,
	}, nil
}
