package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	secret string
}

func NewTokenManager(secret string) *TokenManager {
	return &TokenManager{
		secret: secret,
	}
}

func (tm *TokenManager) Generate(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString([]byte(tm.secret))
}