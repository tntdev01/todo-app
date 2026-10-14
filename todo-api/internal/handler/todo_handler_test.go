package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/tuanntran/todo-app/todo-api/internal/model"
	"github.com/tuanntran/todo-app/todo-api/internal/service"
	"github.com/tuanntran/todo-app/todo-api/internal/testutil"
)

func newTodoHandler(repo *testutil.FakeTodoRepository) *TodoHandler {
	return &TodoHandler{Service: service.NewTodoService(repo)}
}

func performRequest(method, route, path, body string, handlerFunc gin.HandlerFunc) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Handle(method, route, handlerFunc)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, request)
	return recorder
}

func decodeResponse(t *testing.T, recorder *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var response map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response body %q: %v", recorder.Body.String(), err)
	}
	return response
}

func TestTodoHandlerCreate(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		repoErr    error
		wantStatus int
		wantCalls  int
		wantError  string
	}{
		{name: "created", body: `{"title":"Buy milk","comment":"Oat milk"}`, wantStatus: http.StatusCreated, wantCalls: 1},
		{name: "invalid request", body: `{"title":"x","comment":"Oat milk"}`, wantStatus: http.StatusBadRequest},
		{name: "service error", body: `{"title":"Buy milk","comment":"Oat milk"}`, repoErr: errors.New("insert failed"), wantStatus: http.StatusInternalServerError, wantCalls: 1, wantError: "insert failed"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &testutil.FakeTodoRepository{CreateErr: test.repoErr}
			handler := newTodoHandler(repo)
			recorder := performRequest(http.MethodPost, "/todos", "/todos", test.body, handler.Create)

			if recorder.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", recorder.Code, test.wantStatus, recorder.Body.String())
			}
			if repo.CreateCalls != test.wantCalls {
				t.Errorf("Create() repository calls = %d, want %d", repo.CreateCalls, test.wantCalls)
			}
			response := decodeResponse(t, recorder)
			if test.wantError != "" {
				if response["success"] != false || response["error"] != test.wantError {
					t.Errorf("error response = %#v, want success=false and error=%q", response, test.wantError)
				}
			} else if test.wantStatus == http.StatusCreated {
				if response["success"] != true || response["data"] != "created" {
					t.Errorf("success response = %#v, want success=true and data=created", response)
				}
				if repo.CreateTodo.Title != "Buy milk" || repo.CreateTodo.Comment != "Oat milk" {
					t.Errorf("created todo = %+v, want submitted values", repo.CreateTodo)
				}
			}
		})
	}
}

func TestTodoHandlerGetAll(t *testing.T) {
	wantTodos := []model.Todo{{ID: "todo-1", Title: "First"}, {ID: "todo-2", Title: "Second"}}
	tests := []struct {
		name       string
		repoErr    error
		wantStatus int
		wantError  string
	}{
		{name: "success", wantStatus: http.StatusOK},
		{name: "service error", repoErr: errors.New("list failed"), wantStatus: http.StatusInternalServerError, wantError: "list failed"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &testutil.FakeTodoRepository{AllTodos: wantTodos, GetAllErr: test.repoErr}
			handler := newTodoHandler(repo)
			recorder := performRequest(http.MethodGet, "/todos", "/todos", "", handler.GetAll)
			if recorder.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", recorder.Code, test.wantStatus, recorder.Body.String())
			}
			response := decodeResponse(t, recorder)
			if test.wantError != "" {
				if response["success"] != false || response["error"] != test.wantError {
					t.Errorf("error response = %#v, want success=false and error=%q", response, test.wantError)
				}
				return
			}
			if response["success"] != true || repo.GetAllCalls != 1 {
				t.Errorf("success response/calls = %#v/%d, want success=true/1 call", response, repo.GetAllCalls)
			}
			data, err := json.Marshal(response["data"])
			if err != nil {
				t.Fatal(err)
			}
			var gotTodos []model.Todo
			if err := json.Unmarshal(data, &gotTodos); err != nil {
				t.Fatal(err)
			}
			if len(gotTodos) != len(wantTodos) || gotTodos[0].ID != wantTodos[0].ID || gotTodos[1].ID != wantTodos[1].ID {
				t.Errorf("response todos = %+v, want %+v", gotTodos, wantTodos)
			}
		})
	}
}

func TestTodoHandlerGetByID(t *testing.T) {
	wantTodo := model.Todo{ID: "todo-1", Title: "Read"}
	tests := []struct {
		name       string
		repoErr    error
		wantStatus int
		wantError  string
	}{
		{name: "success", wantStatus: http.StatusOK},
		{name: "service error", repoErr: errors.New("not found"), wantStatus: http.StatusInternalServerError, wantError: "not found"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &testutil.FakeTodoRepository{GetByIDTodo: wantTodo, GetByIDErr: test.repoErr}
			handler := newTodoHandler(repo)
			recorder := performRequest(http.MethodGet, "/todos/:id", "/todos/todo-1", "", handler.GetByID)
			if recorder.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", recorder.Code, test.wantStatus, recorder.Body.String())
			}
			response := decodeResponse(t, recorder)
			if test.wantError != "" {
				if response["success"] != false || response["error"] != test.wantError {
					t.Errorf("error response = %#v, want success=false and error=%q", response, test.wantError)
				}
				return
			}
			if repo.RequestedID != "todo-1" || repo.GetByIDCalls != 1 {
				t.Errorf("GetByID() repository call = (%q, %d calls), want (todo-1, 1 call)", repo.RequestedID, repo.GetByIDCalls)
			}
			data, err := json.Marshal(response["data"])
			if err != nil {
				t.Fatal(err)
			}
			var gotTodo model.Todo
			if err := json.Unmarshal(data, &gotTodo); err != nil {
				t.Fatal(err)
			}
			if gotTodo.ID != wantTodo.ID || gotTodo.Title != wantTodo.Title {
				t.Errorf("response todo = %+v, want %+v", gotTodo, wantTodo)
			}
		})
	}
}

func TestTodoHandlerUpdate(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		getErr     error
		updateErr  error
		wantStatus int
		wantError  string
	}{
		{name: "success", body: `{"title":"Updated","comment":"New note","completed":true}`, wantStatus: http.StatusOK},
		{name: "invalid request", body: `{"title":"x","comment":"New note"}`, wantStatus: http.StatusBadRequest},
		{name: "lookup error", body: `{"title":"Updated","comment":"New note"}`, getErr: errors.New("lookup failed"), wantStatus: http.StatusInternalServerError, wantError: "lookup failed"},
		{name: "update error", body: `{"title":"Updated","comment":"New note"}`, updateErr: errors.New("update failed"), wantStatus: http.StatusInternalServerError, wantError: "update failed"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &testutil.FakeTodoRepository{
				GetByIDTodo: model.Todo{ID: "todo-1", Title: "Original", Comment: "Original note"},
				GetByIDErr:  test.getErr,
				UpdateErr:   test.updateErr,
			}
			handler := newTodoHandler(repo)
			recorder := performRequest(http.MethodPut, "/todos/:id", "/todos/todo-1", test.body, handler.Update)
			if recorder.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", recorder.Code, test.wantStatus, recorder.Body.String())
			}
			response := decodeResponse(t, recorder)
			if test.wantError != "" {
				if response["success"] != false || response["error"] != test.wantError {
					t.Errorf("error response = %#v, want success=false and error=%q", response, test.wantError)
				}
				return
			}
			if test.wantStatus == http.StatusBadRequest {
				if repo.GetByIDCalls != 0 || repo.UpdateCalls != 0 {
					t.Errorf("invalid update called repository: lookup=%d update=%d", repo.GetByIDCalls, repo.UpdateCalls)
				}
				return
			}
			if response["success"] != true || response["data"] != "updated" {
				t.Errorf("success response = %#v, want success=true and data=updated", response)
			}
			if repo.RequestedID != "todo-1" || repo.UpdateCalls != 1 {
				t.Errorf("Update() repository calls = lookup(%q), update(%d), want lookup(todo-1), update(1)", repo.RequestedID, repo.UpdateCalls)
			}
			if repo.UpdatedTodo.Title != "Updated" || repo.UpdatedTodo.Comment != "New note" || !repo.UpdatedTodo.Completed {
				t.Errorf("updated todo = %+v, want submitted fields", repo.UpdatedTodo)
			}
		})
	}
}

func TestTodoHandlerDelete(t *testing.T) {
	tests := []struct {
		name       string
		repoErr    error
		wantStatus int
		wantError  string
	}{
		{name: "success", wantStatus: http.StatusOK},
		{name: "service error", repoErr: errors.New("delete failed"), wantStatus: http.StatusInternalServerError, wantError: "delete failed"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &testutil.FakeTodoRepository{DeleteErr: test.repoErr}
			handler := newTodoHandler(repo)
			recorder := performRequest(http.MethodDelete, "/todos/:id", "/todos/todo-1", "", handler.Delete)
			if recorder.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", recorder.Code, test.wantStatus, recorder.Body.String())
			}
			response := decodeResponse(t, recorder)
			if test.wantError != "" {
				if response["success"] != false || response["error"] != test.wantError {
					t.Errorf("error response = %#v, want success=false and error=%q", response, test.wantError)
				}
				return
			}
			if response["success"] != true || response["data"] != "deleted" {
				t.Errorf("success response = %#v, want success=true and data=deleted", response)
			}
			if repo.DeletedID != "todo-1" || repo.DeleteCalls != 1 {
				t.Errorf("Delete() repository call = (%q, %d calls), want (todo-1, 1 call)", repo.DeletedID, repo.DeleteCalls)
			}
		})
	}
}
