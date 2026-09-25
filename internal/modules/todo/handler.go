package todo

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/modasser-nayem/tiny-task/internal/middleware"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateTodoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H {
				"error": "invalid request body",
			},
		)

		return
	}

	userID, ok := middleware.GetUserID(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H {
			"error": "unauthorized",
		})

		return
	}

	todo, err := h.service.Create(c.Request.Context(), userID, req)


	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

func (h *Handler) GetAll(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H {
			"error": "unauthorized",
		})

		return
	}

	todos, err := h.service.GetAll(c.Request.Context(), userID)


	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H {
				"error": "failed to fetch todos",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H {
			"data": todos,
		},
	)
	
}

func (h *Handler) GetByID(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H {
				"error": "invalid todo ID",
			},
		)
		return
	}

	userID, ok := middleware.GetUserID(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	todo, err := h.service.GetByID(c.Request.Context(), userID, id)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H {
				"error": "todo not found",
			},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H {
		"data": todo,
	})
	
}