package repository

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tuanntran/todo-app/todo-api/internal/model"
)

func newTestTodoRepository(t *testing.T) *TodoRepository {
	t.Helper()
	return NewTodoRepository(testPool)
}

func testTodo(title, comment string) model.Todo {
	return model.Todo{
		ID:      uuid.NewString(),
		Title:   title,
		Comment: comment,
	}
}

func cleanupTestTodos(t *testing.T, repo *TodoRepository, todos ...model.Todo) {
	t.Helper()
	t.Cleanup(func() {
		for _, todo := range todos {
			if err := repo.Delete(todo.ID); err != nil {
				t.Errorf("clean up todo %q: %v", todo.ID, err)
			}
		}
	})
}

func TestTodoRepositoryCreate(t *testing.T) {
	repo := newTestTodoRepository(t)
	todo := testTodo("Create integration test", "created by repository test")
	cleanupTestTodos(t, repo, todo)

	if err := repo.Create(todo); err != nil {
		t.Fatalf("Create() error = %v, want nil", err)
	}

	got, err := repo.GetByID(todo.ID)
	if err != nil {
		t.Fatalf("GetByID() after Create() error = %v, want nil", err)
	}
	if got.ID != todo.ID || got.Title != todo.Title || got.Comment != todo.Comment || got.Completed {
		t.Errorf("created todo = %+v, want ID/title/comment from input and Completed=false", got)
	}
	if got.CreatedAt.IsZero() || got.UpdatedAt.IsZero() {
		t.Errorf("database timestamps were not populated: created_at=%v updated_at=%v", got.CreatedAt, got.UpdatedAt)
	}
}

func TestTodoRepositoryGetAll(t *testing.T) {
	repo := newTestTodoRepository(t)
	want := []model.Todo{
		testTodo("First list integration test", "first"),
		testTodo("Second list integration test", "second"),
	}
	cleanupTestTodos(t, repo, want...)
	for _, todo := range want {
		if err := repo.Create(todo); err != nil {
			t.Fatalf("Create(%q) error = %v", todo.ID, err)
		}
	}

	got, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll() error = %v, want nil", err)
	}
	gotByID := make(map[string]model.Todo, len(got))
	for _, todo := range got {
		gotByID[todo.ID] = todo
	}
	for _, todo := range want {
		if gotTodo, ok := gotByID[todo.ID]; !ok || gotTodo.Title != todo.Title || gotTodo.Comment != todo.Comment {
			t.Errorf("GetAll() missing or incorrect todo %q: got %+v, found=%t", todo.ID, gotTodo, ok)
		}
	}
}

func TestTodoRepositoryGetByID(t *testing.T) {
	repo := newTestTodoRepository(t)
	todo := testTodo("Get by ID integration test", "lookup")
	cleanupTestTodos(t, repo, todo)
	if err := repo.Create(todo); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	got, err := repo.GetByID(todo.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v, want nil", err)
	}
	if got.ID != todo.ID || got.Title != todo.Title || got.Comment != todo.Comment {
		t.Errorf("GetByID() = %+v, want todo %+v", got, todo)
	}
}

func TestTodoRepositoryGetByIDMissing(t *testing.T) {
	repo := newTestTodoRepository(t)
	_, err := repo.GetByID(uuid.NewString())
	if !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("GetByID() for missing ID error = %v, want %v", err, pgx.ErrNoRows)
	}
}
func TestTodoRepositoryUpdate(t *testing.T) {
	repo := newTestTodoRepository(t)
	todo := testTodo("Original integration test title", "original")
	cleanupTestTodos(t, repo, todo)
	if err := repo.Create(todo); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	originalUpdatedAt, err := repo.GetByID(todo.ID)
	if err != nil {
		t.Fatalf("GetByID() before Update() error = %v", err)
	}
	todo.Title = "Updated integration test title"
	todo.Comment = "updated"
	todo.Completed = true

	if err := repo.Update(todo); err != nil {
		t.Fatalf("Update() error = %v, want nil", err)
	}
	got, err := repo.GetByID(todo.ID)
	if err != nil {
		t.Fatalf("GetByID() after Update() error = %v", err)
	}
	if got.Title != todo.Title || got.Comment != todo.Comment || got.Completed != todo.Completed {
		t.Errorf("updated todo = %+v, want title/comment/completed from input", got)
	}
	if got.UpdatedAt.Before(originalUpdatedAt.UpdatedAt) {
		t.Errorf("updated_at moved backwards: before=%v after=%v", originalUpdatedAt.UpdatedAt, got.UpdatedAt)
	}
}

func TestTodoRepositoryDelete(t *testing.T) {
	repo := newTestTodoRepository(t)
	todo := testTodo("Delete integration test", "delete me")
	if err := repo.Create(todo); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := repo.Delete(todo.ID); err != nil {
		t.Fatalf("Delete() error = %v, want nil", err)
	}
	if _, err := repo.GetByID(todo.ID); !errors.Is(err, pgx.ErrNoRows) {
		t.Errorf("GetByID() after Delete() error = %v, want %v", err, pgx.ErrNoRows)
	}
}

func TestTodoRepositoryDeleteMissing(t *testing.T) {
	repo := newTestTodoRepository(t)
	err := repo.Delete(uuid.NewString())
	if err != nil {
		t.Fatalf("Delete() for missing ID error = %v, want nil", err)
	}
}

func TestTodoRepositoryCreateRejectsDuplicateID(t *testing.T) {
	repo := newTestTodoRepository(t)
	todo := testTodo("Duplicate integration test", "duplicate ID")
	cleanupTestTodos(t, repo, todo)
	if err := repo.Create(todo); err != nil {
		t.Fatalf("first Create() error = %v, want nil", err)
	}

	duplicate := todo
	duplicate.Title = "Duplicate"
	if err := repo.Create(duplicate); err == nil {
		t.Fatal("second Create() error = nil, want unique ID constraint error")
	}
}

func TestTodoRepositoryUpdateMissingID(t *testing.T) {
	repo := newTestTodoRepository(t)
	todo := model.Todo{ID: uuid.NewString(), Title: "Missing integration test", Completed: true}

	if err := repo.Update(todo); err != nil {
		t.Fatalf("Update() error = %v, want nil for a missing ID", err)
	}
}

func TestTodoRepositoryTimestampsUseDatabaseTime(t *testing.T) {
	repo := newTestTodoRepository(t)
	todo := testTodo("Timestamp integration test", "time check")
	cleanupTestTodos(t, repo, todo)

	beforeCreate := time.Now().UTC().Add(-time.Minute)
	if err := repo.Create(todo); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	got, err := repo.GetByID(todo.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if got.CreatedAt.Before(beforeCreate) || got.UpdatedAt.Before(beforeCreate) {
		t.Errorf("database timestamps precede creation window: created_at=%v updated_at=%v", got.CreatedAt, got.UpdatedAt)
	}
}
