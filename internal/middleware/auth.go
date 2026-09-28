package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	apperrors "github.com/modasser-nayem/tiny-task/internal/errors"
	"github.com/modasser-nayem/tiny-task/internal/modules/auth"
)

const UserIDKey = "userID"

func Auth(tokenManager *auth.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.Error(apperrors.Unauthorized(
				"AUTHORIZATION_HEADER_REQUIRED",
				"authorization header is required",
			))

			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Error(apperrors.Unauthorized(
				"INVALID_AUTHORIZATION_HEADER",
				"invalid authorization header",
			))

			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := tokenManager.Parse(tokenString)

		if err != nil {
			c.Error(apperrors.Unauthorized(
				"INVALID_TOKEN",
				"invalid or expired token",
			))

			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Next()
	}
}



func GetUserID(c *gin.Context) (int64, bool) {
	value, exists := c.Get(UserIDKey)

	if !exists {
		return 0, false
	}

	userID, ok := value.(int64)

	return userID, ok
}