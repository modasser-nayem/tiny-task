package todo

import (
	"strings"

	apperrors "github.com/modasser-nayem/tiny-task/internal/errors"
)

func validateCreateTodo(req CreateTodoRequest) error {

	title := strings.TrimSpace(req.Title)

	if title == "" {
		return apperrors.BadRequest(
			"TITLE_REQUIRED",
			"title is required",
		)
	}

	if len(title) > 255 {
		return apperrors.BadRequest(
			"TITLE_TOO_LONG",
			"title cannot exceed 255 characters",
		)
	}

	return nil
}
