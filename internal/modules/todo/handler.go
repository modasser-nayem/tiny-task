package todo

import (
	"errors"
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

	query, err := parseListQuery(c)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	todos, err := h.service.GetAll(c.Request.Context(), userID, query)


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
			"meta": gin.H {
				"total": len(todos),
				"page": query.Page,
				"limit": query.Limit,
			},
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

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(
		c.Param("id"),
		10,
		64,
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid todo ID",
			},
		)
		return
	}

	var req UpdateTodoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid request body",
			},
		)
		return
	}

	userID, ok := middleware.GetUserID(c)

	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "unauthorized",
			},
		)
		return
	}

	todo, err := h.service.Update(
		c.Request.Context(),
		userID,
		id,
		req,
	)

	if err != nil {
		if errors.Is(err, ErrTodoNotFound) {
			c.JSON(
				http.StatusNotFound,
				gin.H{
					"error": "todo not found",
				},
			)
			return
		}

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		todo,
	)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(
		c.Param("id"),
		10,
		64,
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid todo ID",
			},
		)
		return
	}

	userID, ok := middleware.GetUserID(c)

	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "unauthorized",
			},
		)
		return
	}

	err = h.service.Delete(
		c.Request.Context(),
		userID,
		id,
	)

	if err != nil {
		if errors.Is(err, ErrTodoNotFound) {
			c.JSON(
				http.StatusNotFound,
				gin.H{
					"error": "todo not found",
				},
			)
			return
		}

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "failed to delete todo",
			},
		)
		return
	}

	c.Status(http.StatusNoContent)
}


// parse query parameters
func parseListQuery(c *gin.Context) (ListTodoQuery, error) {
	page := 1
	limit := 10

	if value := c.Query("page"); value != "" {
		parsed, err := strconv.Atoi(value)

		if err != nil || parsed < 0 {
			return ListTodoQuery{}, errors.New("page must be a positive integer")
		}

		page = parsed
	}

	if value := c.Query("limit"); value != "" {
		parsed, err := strconv.Atoi(value)

		if err != nil || parsed <= 0 {
			return ListTodoQuery{}, errors.New("limit must be a positive integer")
		}

		if parsed > 100 {
			return ListTodoQuery{}, errors.New("limit cannot exceed 100")
		}

		limit = parsed
	}

	var completed *bool
	if value := c.Query("completed"); value != "" {
		parsed, err := strconv.ParseBool(value)

		if err != nil {
			return ListTodoQuery{}, errors.New("completed must be a boolean")
		}

		completed = &parsed
	}

	return ListTodoQuery{
		Page: page,
		Limit: limit,
		Completed: completed,
	}, nil
}