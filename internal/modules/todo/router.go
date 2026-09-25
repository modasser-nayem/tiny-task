package todo

import (
	"github.com/gin-gonic/gin"
	"github.com/modasser-nayem/tiny-task/internal/middleware"
	"github.com/modasser-nayem/tiny-task/internal/modules/auth"
)

func RegisterRoutes(r *gin.RouterGroup, handler *Handler, tokenManager *auth.TokenManager) {
	todoRoutes := r.Group("/todos")

	todoRoutes.Use(middleware.Auth(tokenManager))


	todoRoutes.POST("/", handler.Create)
	todoRoutes.GET("/", handler.GetAll)
	todoRoutes.GET("/:id", handler.GetByID)
	// todoRoutes.PUT("/:id", handler.Update)
	// todoRoutes.DELETE("/:id", handler.Delete)

}