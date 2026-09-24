package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/modasser-nayem/tiny-task/internal/modules/auth"
)

func AuthMiddleware(tokenManager *auth.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		
	}
}