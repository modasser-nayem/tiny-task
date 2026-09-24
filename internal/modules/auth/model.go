package auth

import "time"

type User struct {
	ID           int64
	Name         string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}