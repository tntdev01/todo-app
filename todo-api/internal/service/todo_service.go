package service

import (
	"github.com/bytedance/gopkg/util/logger"
	"github.com/google/uuid"

	"github.com/tuanntran/todo-app/todo-api/internal/model"
)

type TodoRepository interface {
	Create(todo model.Todo) error
	GetAll() ([]model.Todo, error)
	GetByID(id string) (model.Todo, error)
	Update(todo model.Todo) error
	Delete(id string) error
}

type TodoService struct {
	repo TodoRepository
}

func NewTodoService(repo TodoRepository) *TodoService {
	return &TodoService{
		repo: repo,
	}
}

func (s *TodoService) Create(
	title string,
	comment string,
) error {

	todo := model.Todo{
		ID:        uuid.New().String(),
		Title:     title,
		Comment:   comment,
		Completed: false,
	}

	err := s.repo.Create(todo)
	if err != nil {
		return err
	}

	logger.Info(
		"todo created",
		"todo_id", todo.ID,
		"title", todo.Title,
		"comment", todo.Comment,
		"completed", todo.Completed,
	)

	return nil
}

func (s *TodoService) GetAll() ([]model.Todo, error) {
	return s.repo.GetAll()
}

func (s *TodoService) GetByID(id string) (model.Todo, error) {
	return s.repo.GetByID(id)
}

func (s *TodoService) Update(
	id string,
	title string,
	comment string,
	completed bool,
) error {
	todo, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	todo.Title = title
	todo.Comment = comment
	todo.Completed = completed

	return s.repo.Update(todo)
}

func (s *TodoService) Delete(id string) error {
	return s.repo.Delete(id)
}
