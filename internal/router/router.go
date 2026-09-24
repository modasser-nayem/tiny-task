package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/modasser-nayem/tiny-task/internal/config"
	"github.com/modasser-nayem/tiny-task/internal/modules/auth"
)

func Setup(cfg config.Config, authHandler *auth.Handler) *gin.Engine {
  
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	api := r.Group("/api/v1")

	auth.RegisterRoutes(api, authHandler)

	return r
}