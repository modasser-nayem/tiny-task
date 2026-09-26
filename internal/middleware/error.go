package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/modasser-nayem/tiny-task/internal/errors"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err


		var appErr *apperrors.AppError

		if errors.As(err, &appErr) {
			c.JSON(
				appErr.Status,
				gin.H {
					"error": gin.H {
						"code": appErr.Code,
						"message": appErr.Message,
					},
				},
			)
			return
		}

		log.Println("Unhandled Error: ", err)
		
		c.JSON(
			http.StatusInternalServerError, 
			gin.H{
			"error": gin.H {
				"code": "internal",
				"message": "Internal server error",
			},
		})
	}
}