package todo

import (
	"context"
	"errors"
	"strings"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(ctx context.Context, userID int64, req CreateTodoRequest) (*TodoResponse, error) {
	title := strings.TrimSpace(req.Title)

	if title == "" {
		return nil, errors.New("title is required")
	}

	todo := &Todo{
		UserID: userID,
		Title: title,
		Description: req.Description,
	}

	created, err := s.repository.Create(ctx, todo)

	if err != nil {
		return nil, err
	}

	return toResponse(created), nil
}

func (s *Service) GetAll(ctx context.Context, userID int64 ) ([]TodoResponse, error) {
	todos, err := s.repository.FindAllByUserID(ctx, userID)

	if err != nil {
		return nil, err
	}

	responses := make([]TodoResponse, 0, len(todos))

	for i := range todos {
		responses = append(responses, *toResponse(&todos[i]))
	}

	return responses, nil
}

func (s *Service) GetByID(ctx context.Context, userID int64, id int64) (*TodoResponse, error) {
	todo, err := s.repository.FindByID(ctx, id, userID)

	if err != nil {
		return nil, err
	}

	return toResponse(todo), nil
}

func toResponse(todo *Todo) *TodoResponse {
	return &TodoResponse{
		ID: todo.ID,
		UserID: todo.UserID,
		Title: todo.Title,
		Description: todo.Description,
		Completed: todo.Completed,
	}
}