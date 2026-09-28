package auth

import (
	"context"
	"errors"
	"testing"
)

type mockAuthRepository struct {
	findByEmailFn func(ctx context.Context, email string) (*User, error)
	createFn      func(ctx context.Context, user *User) (*User, error)
}

func (m *mockAuthRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	return m.findByEmailFn(ctx, email)
}

func (m *mockAuthRepository) Create(ctx context.Context, user *User) (*User, error) {
	return m.createFn(ctx, user)
}

func TestRegisterValidation(t *testing.T) {
	repo := &mockAuthRepository{}
	tokenManager := NewTokenManager("12345678901234567890123456789012")
	service := NewService(repo, tokenManager)

	// Missing name
	_, err := service.Register(context.Background(), RegisterRequest{
		Name:     "   ",
		Email:    "test@example.com",
		Password: "password123",
	})
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}

	// Missing email
	_, err = service.Register(context.Background(), RegisterRequest{
		Name:     "John",
		Email:    "  ",
		Password: "password123",
	})
	if err == nil {
		t.Fatal("expected error for empty email, got nil")
	}

	// Password too short
	_, err = service.Register(context.Background(), RegisterRequest{
		Name:     "John",
		Email:    "john@example.com",
		Password: "123",
	})
	if err == nil {
		t.Fatal("expected error for short password, got nil")
	}
}

func TestRegisterDuplicateEmail(t *testing.T) {
	repo := &mockAuthRepository{
		findByEmailFn: func(ctx context.Context, email string) (*User, error) {
			return &User{ID: 1, Email: email}, nil
		},
	}
	tokenManager := NewTokenManager("12345678901234567890123456789012")
	service := NewService(repo, tokenManager)

	_, err := service.Register(context.Background(), RegisterRequest{
		Name:     "John",
		Email:    "existing@example.com",
		Password: "password123",
	})
	if !errors.Is(err, ErrEmailAlreadyExists) {
		t.Fatalf("expected ErrEmailAlreadyExists, got %v", err)
	}
}

func TestRegisterSuccess(t *testing.T) {
	repo := &mockAuthRepository{
		findByEmailFn: func(ctx context.Context, email string) (*User, error) {
			return nil, ErrUserNotFound
		},
		createFn: func(ctx context.Context, user *User) (*User, error) {
			user.ID = 10
			return user, nil
		},
	}
	tokenManager := NewTokenManager("12345678901234567890123456789012")
	service := NewService(repo, tokenManager)

	res, err := service.Register(context.Background(), RegisterRequest{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "securePassword123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.ID != 10 {
		t.Errorf("expected ID 10, got %d", res.ID)
	}
	if res.Email != "alice@example.com" {
		t.Errorf("expected email alice@example.com, got %s", res.Email)
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	repo := &mockAuthRepository{
		findByEmailFn: func(ctx context.Context, email string) (*User, error) {
			return nil, ErrUserNotFound
		},
	}
	tokenManager := NewTokenManager("12345678901234567890123456789012")
	service := NewService(repo, tokenManager)

	_, err := service.Login(context.Background(), &LoginRequest{
		Email:    "unknown@example.com",
		Password: "password123",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}
