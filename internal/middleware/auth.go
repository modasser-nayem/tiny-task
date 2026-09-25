package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/modasser-nayem/tiny-task/internal/modules/auth"
)

const UserIDKey = "userID"

func Auth(tokenManager *auth.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "authorization header is required",
				},
			)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H {
					"error": "invalid authorization header",
				},
			)
			return
		}

		tokenString := parts[1]

		claims, err := tokenManager.Parse(tokenString)

		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H {
					"error": "invalid or expired token",
				},
			)
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