package user

import (
	"github.com/gin-gonic/gin"
	"github.com/modasser-nayem/tiny-task/internal/middleware"
	"github.com/modasser-nayem/tiny-task/internal/modules/auth"
)

func RegisterRoutes(router *gin.RouterGroup, handler *Handler, tokenManager *auth.TokenManager) {
	userRoutes := router.Group("users")

	userRoutes.Use(middleware.Auth(tokenManager))

	userRoutes.GET("/me", handler.Me)

}