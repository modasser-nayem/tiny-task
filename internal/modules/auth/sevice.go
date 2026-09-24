package auth

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository Repository
	tokenManager *TokenManager
}

func NewService(repository Repository, tokenManager *TokenManager) *Service {
	return &Service{
		repository: repository,
		tokenManager: tokenManager,
	}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	name := strings.TrimSpace(req.Name)
	email := strings.ToLower(strings.TrimSpace(req.Email))

	if name == "" {
		return nil, errors.New("name is required")
	}

	if email == "" {
		return nil, errors.New("email is required")
	}

	if req.Password == "" {
		return nil, errors.New("password is required")
	}

	existingUser, err := s.repository.FindByEmail(ctx, email)

	if err != nil && existingUser != nil {
		return nil, err
	}

	if !errors.Is(err, ErrUserNotFound) && existingUser != nil {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &User{
		Name: name,
		Email: email,
		PasswordHash: string(passwordHash),
	}

	createdUser, err := s.repository.Create(ctx, user)

	if err != nil {
		return nil, err
	}

	return &RegisterResponse{
		ID: createdUser.ID,
		Name: createdUser.Name,
		Email: createdUser.Email,
	}, nil

}

func (s *Service) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	if email == "" || req.Password == "" {
		return nil, errors.New("email and password are required")
	}

	user, err := s.repository.FindByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, errors.New("invalid credentials")
		}

		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))


	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := s.tokenManager.Generate(user.ID)

	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	return &LoginResponse {
		AccessToken: accessToken,
	}, nil
}

