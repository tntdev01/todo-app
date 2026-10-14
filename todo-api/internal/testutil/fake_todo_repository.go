package testutil

import "github.com/tuanntran/todo-app/todo-api/internal/model"

type FakeTodoRepository struct {
	CreateTodo   model.Todo
	CreateErr    error
	AllTodos     []model.Todo
	GetAllErr    error
	GetByIDTodo  model.Todo
	RequestedID  string
	GetByIDErr   error
	UpdatedTodo  model.Todo
	UpdateErr    error
	DeletedID    string
	DeleteErr    error
	CreateCalls  int
	GetAllCalls  int
	GetByIDCalls int
	UpdateCalls  int
	DeleteCalls  int
}

func (r *FakeTodoRepository) Create(todo model.Todo) error {
	r.CreateCalls++
	r.CreateTodo = todo
	return r.CreateErr
}

func (r *FakeTodoRepository) GetAll() ([]model.Todo, error) {
	r.GetAllCalls++
	return r.AllTodos, r.GetAllErr
}

func (r *FakeTodoRepository) GetByID(id string) (model.Todo, error) {
	r.GetByIDCalls++
	r.RequestedID = id
	return r.GetByIDTodo, r.GetByIDErr
}

func (r *FakeTodoRepository) Update(todo model.Todo) error {
	r.UpdateCalls++
	r.UpdatedTodo = todo
	return r.UpdateErr
}

func (r *FakeTodoRepository) Delete(id string) error {
	r.DeleteCalls++
	r.DeletedID = id
	return r.DeleteErr
}
