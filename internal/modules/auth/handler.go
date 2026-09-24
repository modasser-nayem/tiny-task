package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"error": "invalid request body",
		})
		return
	}

	ctx := c.Request.Context()

	response, err := h.service.Register(ctx, req)

	if err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H {
				"error": "email already exists",
			})
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}


func(h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})

		return
	}

	ctx := c.Request.Context()

	response, err := h.service.Login(ctx, &req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, response)
}