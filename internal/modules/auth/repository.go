package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserNotFound = errors.New("user not found")
var ErrEmailAlreadyExists = errors.New("email already exists")

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// FindByEmail retrieves a user by email.
// It returns ErrUserNotFound if the user is not found.
func (r *PostgresRepository) FindByEmail(ctx context.Context, email string) (*User, error) {


	query := `SELECT id, name, email, password_hash, created_at, updated_at FROM users WHERE email = $1`

	var user User

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &user, nil
}

// Create adds a new user to the database.
func(r *PostgresRepository) Create(ctx context.Context, user *User) (*User, error) {
	query := `
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query, user.Name, user.Email, user.PasswordHash).Scan(
		&user.ID, 
		&user.CreatedAt, 
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}