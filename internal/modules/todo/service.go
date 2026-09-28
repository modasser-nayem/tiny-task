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

if err := validateCreateTodo(req); err != nil {
		return nil, err
	}

	req.Title = strings.TrimSpace(req.Title)

	todo := &Todo{
		UserID: userID,
		Title: req.Title,
		Description: req.Description,
	}

	created, err := s.repository.Create(ctx, todo)

	if err != nil {
		return nil, err
	}

	return toResponse(created), nil
}

func (s *Service) GetAll(ctx context.Context, userID int64, query ListTodoQuery ) ([]TodoResponse, error) {
	todos, err := s.repository.FindAllByUserID(ctx, userID, query)

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

func (s *Service) Update(ctx context.Context, userID int64, id int64, req UpdateTodoRequest) (*TodoResponse, error) {
	todo, err := s.repository.FindByID(ctx, id, userID)

	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)
		if title == "" {
			return nil, errors.New("title cannot be empty")
		}
		todo.Title = title
	}


	if req.Description != nil {
		todo.Description = req.Description
	}

	if req.Completed != nil {
		todo.Completed = *req.Completed
	}

	updated, err := s.repository.Update(ctx, todo)
  if err != nil {
		return nil, err
	}

	return toResponse(updated), nil
}

func (s *Service) Delete(
	ctx context.Context,
	userID int64,
	id int64,
) error {

	return s.repository.Delete(
		ctx,
		id,
		userID,
	)
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