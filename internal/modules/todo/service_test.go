package todo

import (
	"context"
	"errors"
	"testing"
)

type mockRepository struct {
    createFn  func(ctx context.Context, todo *Todo) (*Todo, error)
    findAllFn func(ctx context.Context, userID int64, query ListTodoQuery) ([]Todo, error)
    findByIDFn func(ctx context.Context, id int64, userID int64) (*Todo, error)
    updateFn  func(ctx context.Context, todo *Todo) (*Todo, error)
    deleteFn  func(ctx context.Context, id int64, userID int64) error
}
var _ Repository = (*mockRepository)(nil)

func (m *mockRepository) Create(
    ctx context.Context,
    todo *Todo,
) (*Todo, error) {
    return m.createFn(ctx, todo)
}

func (m *mockRepository) FindAllByUserID(
    ctx context.Context,
    userID int64,
    query ListTodoQuery,
) ([]Todo, error) {
    return m.findAllFn(ctx, userID, query)
}

func (m *mockRepository) FindByID(
    ctx context.Context,
    id int64,
    userID int64,
) (*Todo, error) {
    return m.findByIDFn(ctx, id, userID)
}

func (m *mockRepository) Update(
    ctx context.Context,
    todo *Todo,
) (*Todo, error) {
    return m.updateFn(ctx, todo)
}

func (m *mockRepository) Delete(
    ctx context.Context,
    id int64,
    userID int64,
) error {
    return m.deleteFn(ctx, id, userID)
}

var errMockRepository = errors.New("mock repository error")

func TestCreateRejectsEmptyTitle(t *testing.T) {
    repo := &mockRepository{
        createFn: func(
            ctx context.Context,
            todo *Todo,
        ) (*Todo, error) {
            t.Fatal("repository should not be called")
            return nil, nil
        },
    }

    service := NewService(repo)

    _, err := service.Create(
        context.Background(),
        1,
        CreateTodoRequest{
            Title: "   ",
        },
    )

    if err == nil {
        t.Fatal("expected validation error, got nil")
    }
}

func TestCreateTodoSuccess(t *testing.T) {
    repo := &mockRepository{
        createFn: func(
            ctx context.Context,
            todo *Todo,
        ) (*Todo, error) {
            todo.ID = 101
            return todo, nil
        },
    }

    service := NewService(repo)

    result, err := service.Create(
        context.Background(),
        5,
        CreateTodoRequest{
            Title: "  Learn Go  ",
        },
    )

    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }

    if result.ID != 101 {
        t.Errorf("expected ID 101, got %d", result.ID)
    }

    if result.Title != "Learn Go" {
        t.Errorf("expected trimmed title, got %q", result.Title)
    }

    if result.UserID != 5 {
        t.Errorf("expected user ID 5, got %d", result.UserID)
    }
}


func TestCreateTodoRepositoryError(t *testing.T) {
    repo := &mockRepository{
        createFn: func(
            ctx context.Context,
            todo *Todo,
        ) (*Todo, error) {
            return nil, errMockRepository
        },
    }

    service := NewService(repo)

    _, err := service.Create(
        context.Background(),
        5,
        CreateTodoRequest{
            Title: "Learn Go",
        },
    )

    if err == nil {
        t.Fatal("expected repository error, got nil")
    }
}