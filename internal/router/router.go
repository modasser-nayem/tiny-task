package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/modasser-nayem/tiny-task/internal/config"
	"github.com/modasser-nayem/tiny-task/internal/modules/auth"
	"github.com/modasser-nayem/tiny-task/internal/modules/todo"
	"github.com/modasser-nayem/tiny-task/internal/modules/user"
)

func Setup(cfg config.Config, authHandler *auth.Handler, userHandler *user.Handler, todoHandler *todo.Handler, tokenManager *auth.TokenManager) *gin.Engine {
  
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	api := r.Group("/api/v1")

	auth.RegisterRoutes(api, authHandler)
	user.RegisterRoutes(api, userHandler, tokenManager)
	todo.RegisterRoutes(api, todoHandler, tokenManager)

	return r
}