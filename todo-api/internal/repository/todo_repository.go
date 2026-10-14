package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tuanntran/todo-app/todo-api/internal/model"
)

type TodoRepository struct {
	DB *pgxpool.Pool
}

func NewTodoRepository(db *pgxpool.Pool) *TodoRepository {
	return &TodoRepository{
		DB: db,
	}
}

func (r *TodoRepository) Create(todo model.Todo) error {
	_, err := r.DB.Exec(
		context.Background(),
		`INSERT INTO todos
        (id,title,comment,completed)
        VALUES($1,$2,$3,$4)`,
		todo.ID,
		todo.Title,
		todo.Comment,
		todo.Completed,
	)

	return err
}

func (r *TodoRepository) GetAll() ([]model.Todo, error) {

	rows, err := r.DB.Query(
		context.Background(),
		`SELECT id,title,comment,completed,
        created_at,updated_at
        FROM todos`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	todos := make([]model.Todo, 0)

	for rows.Next() {

		var t model.Todo

		err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Comment,
			&t.Completed,
			&t.CreatedAt,
			&t.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		todos = append(todos, t)
	}

	return todos, nil
}

func (r *TodoRepository) GetByID(id string) (model.Todo, error) {

	var t model.Todo

	err := r.DB.QueryRow(
		context.Background(),
		`SELECT id,title,comment,completed,
         created_at,updated_at
        FROM todos
         WHERE id=$1`,
		id,
	).Scan(
		&t.ID,
		&t.Title,
		&t.Comment,
		&t.Completed,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	return t, err
}

func (r *TodoRepository) Update(todo model.Todo) error {
	_, err := r.DB.Exec(
		context.Background(),
		`UPDATE todos
        SET title=$1, comment=$2, completed=$3, updated_at=NOW()
        WHERE id=$4`,
		todo.Title,
		todo.Comment,
		todo.Completed,
		todo.ID,
	)

	return err
}

func (r *TodoRepository) Delete(id string) error {
	_, err := r.DB.Exec(
		context.Background(),
		`DELETE FROM todos WHERE id=$1`,
		id,
	)

	return err
}
