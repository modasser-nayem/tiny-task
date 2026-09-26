package todo

type CreateTodoRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

type UpdateTodoRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}

type TodoResponse struct {
	ID          int64   `json:"id"`
	UserID      int64   `json:"user_id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Completed   bool    `json:"completed"`
}

type ListTodoQuery struct {
	Page int
	Limit int
	Completed *bool
}