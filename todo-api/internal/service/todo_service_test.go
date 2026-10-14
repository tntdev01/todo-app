package service

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/tuanntran/todo-app/todo-api/internal/model"
	"github.com/tuanntran/todo-app/todo-api/internal/testutil"
)

func TestTodoServiceCreate(t *testing.T) {
	wantErr := errors.New("create failed")
	repo := &testutil.FakeTodoRepository{CreateErr: wantErr}
	service := NewTodoService(repo)

	err := service.Create("Buy milk", "Oat milk")

	if err != wantErr {
		t.Fatalf("Create() error = %v, want %v", err, wantErr)
	}
	if repo.CreateCalls != 1 {
		t.Fatalf("Create() called repository %d times, want 1", repo.CreateCalls)
	}
	if _, err := uuid.Parse(repo.CreateTodo.ID); err != nil {
		t.Fatalf("Create() generated invalid UUID %q: %v", repo.CreateTodo.ID, err)
	}
	if repo.CreateTodo.Title != "Buy milk" || repo.CreateTodo.Comment != "Oat milk" {
		t.Fatalf("Create() stored todo = %+v, want supplied title and comment", repo.CreateTodo)
	}
	if repo.CreateTodo.Completed {
		t.Fatal("Create() stored todo as completed, want false")
	}
}

func TestTodoServiceGetAll(t *testing.T) {
	wantTodos := []model.Todo{{ID: "one", Title: "First"}, {ID: "two", Title: "Second"}}
	wantErr := errors.New("list failed")
	repo := &testutil.FakeTodoRepository{AllTodos: wantTodos, GetAllErr: wantErr}
	service := NewTodoService(repo)

	gotTodos, err := service.GetAll()

	if !reflect.DeepEqual(gotTodos, wantTodos) {
		t.Errorf("GetAll() todos = %+v, want %+v", gotTodos, wantTodos)
	}
	if err != wantErr {
		t.Errorf("GetAll() error = %v, want %v", err, wantErr)
	}
	if repo.GetAllCalls != 1 {
		t.Errorf("GetAll() called repository %d times, want 1", repo.GetAllCalls)
	}
}

func TestTodoServiceGetByID(t *testing.T) {
	wantTodo := model.Todo{ID: "todo-1", Title: "Read"}
	wantErr := errors.New("not found")
	repo := &testutil.FakeTodoRepository{GetByIDTodo: wantTodo, GetByIDErr: wantErr}
	service := NewTodoService(repo)

	gotTodo, err := service.GetByID("todo-1")

	if !reflect.DeepEqual(gotTodo, wantTodo) {
		t.Errorf("GetByID() todo = %+v, want %+v", gotTodo, wantTodo)
	}
	if err != wantErr {
		t.Errorf("GetByID() error = %v, want %v", err, wantErr)
	}
	if repo.RequestedID != "todo-1" || repo.GetByIDCalls != 1 {
		t.Errorf("GetByID() repository call = (%q, %d calls), want (%q, 1 call)", repo.RequestedID, repo.GetByIDCalls, "todo-1")
	}
}

func TestTodoServiceUpdate(t *testing.T) {
	createdAt := time.Date(2025, time.January, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	existingTodo := model.Todo{
		ID:        "todo-1",
		Title:     "Old title",
		Comment:   "Old comment",
		Completed: false,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	wantErr := errors.New("update failed")
	repo := &testutil.FakeTodoRepository{GetByIDTodo: existingTodo, UpdateErr: wantErr}
	service := NewTodoService(repo)

	err := service.Update("todo-1", "New title", "New comment", true)

	if err != wantErr {
		t.Errorf("Update() error = %v, want %v", err, wantErr)
	}
	if repo.RequestedID != "todo-1" || repo.GetByIDCalls != 1 {
		t.Errorf("Update() lookup = (%q, %d calls), want (%q, 1 call)", repo.RequestedID, repo.GetByIDCalls, "todo-1")
	}
	wantTodo := existingTodo
	wantTodo.Title = "New title"
	wantTodo.Comment = "New comment"
	wantTodo.Completed = true
	if !reflect.DeepEqual(repo.UpdatedTodo, wantTodo) {
		t.Errorf("Update() saved todo = %+v, want %+v", repo.UpdatedTodo, wantTodo)
	}
	if repo.UpdateCalls != 1 {
		t.Errorf("Update() called repository update %d times, want 1", repo.UpdateCalls)
	}
}

func TestTodoServiceUpdateStopsWhenLookupFails(t *testing.T) {
	wantErr := errors.New("lookup failed")
	repo := &testutil.FakeTodoRepository{GetByIDErr: wantErr}
	service := NewTodoService(repo)

	err := service.Update("todo-1", "New title", "New comment", true)

	if err != wantErr {
		t.Errorf("Update() error = %v, want %v", err, wantErr)
	}
	if repo.UpdateCalls != 0 {
		t.Errorf("Update() called repository update %d times after lookup failure, want 0", repo.UpdateCalls)
	}
}

func TestTodoServiceDelete(t *testing.T) {
	wantErr := errors.New("delete failed")
	repo := &testutil.FakeTodoRepository{DeleteErr: wantErr}
	service := NewTodoService(repo)

	err := service.Delete("todo-1")

	if err != wantErr {
		t.Errorf("Delete() error = %v, want %v", err, wantErr)
	}
	if repo.DeletedID != "todo-1" || repo.DeleteCalls != 1 {
		t.Errorf("Delete() repository call = (%q, %d calls), want (%q, 1 call)", repo.DeletedID, repo.DeleteCalls, "todo-1")
	}
}
