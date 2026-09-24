package todo

type CreateTodoDto struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}