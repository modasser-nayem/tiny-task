package todo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrTodoNotFound = errors.New("todo not found")

type Repository interface {
	Create(ctx context.Context, todo *Todo) (*Todo, error)

	FindAllByUserID(ctx context.Context, userID int64) ([]Todo, error)

	FindByID(ctx context.Context, id int64, userID int64) (*Todo, error)

	Update(ctx context.Context, todo *Todo) (*Todo, error)

	Delete(ctx context.Context, id int64, userID int64) (error)
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Create 
func(r *PostgresRepository) Create(ctx context.Context, todo *Todo) (*Todo, error) {
	var created Todo

	query := `
	INSERT INTO todos (user_id, title, description, completed)
	VALUES ($1, $2, $3, $4)
	RETURNING id, user_id, title, description, completed, created_at, updated_at
	`
err := r.db.QueryRow(ctx, query, todo.UserID, todo.Title, todo.Description, todo.Completed).Scan(
	&created.ID,
	&created.UserID,
	&created.Title,
	&created.Description,
	&created.Completed,
	&created.CreatedAt,
	&created.UpdatedAt,
)

if err != nil {
	return nil, err
}

return &created, nil

}

func (r *PostgresRepository) FindByID(ctx context.Context, id int64, userID int64) (*Todo, error) {
	var todo Todo

	query := `
	SELECT id, user_id, title, description, completed, created_at, updated_at
	FROM todos
	WHERE id = $1
	AND user_id = $2
	`

	err := r.db.QueryRow(ctx, query, id, userID).Scan(
		&todo.ID,
		&todo.UserID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTodoNotFound
	}

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func(r *PostgresRepository) FindAllByUserID(ctx context.Context, userID int64) ([]Todo, error) {

	query := `
	SELECT id, user_id, title, description, completed, created_at, updated_at
	FROM todos
	WHERE user_id = $1
	ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	todos := make([]Todo, 0)

	for rows.Next() {
		var todo Todo

		err := rows.Scan(
			&todo.ID,
			&todo.UserID,
			&todo.Title,
			&todo.Description,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
			return nil, err
		}

		return todos, nil
}

func (r *PostgresRepository) Update(ctx context.Context, todo *Todo) (*Todo, error) {
	var updated Todo

	query := 
	`UPDATE todos
	 SET title = $1, description = $2, completed = $3, updated_at = NOW()
	 WHERE id = $4
	 AND user_id = $5
	 RETURNING id, user_id, title, description, completed, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query, todo.Title, todo.Description, todo.Completed, todo.ID, todo.UserID).Scan(
		&updated.ID,
		&updated.UserID,
		&updated.Title,
		&updated.Description,
		&updated.Completed,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func(r *PostgresRepository) Delete(ctx context.Context, id int64, userID int64) error {
	result, err := r.db.Exec(ctx, `
	DELETE FROM todos
	WHERE id = $1
	AND user_id = $2`, id, userID)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrTodoNotFound
	}

	return nil
}

