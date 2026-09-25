package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/modasser-nayem/tiny-task/internal/middleware"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Me(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)

	if !exists {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "user ID not found in context",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"user_id": userID,
		},
	)
}