package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Nikita527/todo/internal/api"
	"github.com/Nikita527/todo/internal/todo"
)

func TestAPI(t *testing.T) {
	store := todo.New()
	file := filepath.Join(t.TempDir(), "todos.json")
	handler := api.New(store, file).Routes()

	t.Run("list empty", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}

		var tasks []todo.Task
		if err := json.NewDecoder(rec.Body).Decode(&tasks); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if len(tasks) != 0 {
			t.Fatalf("len(tasks) = %d, want 0", len(tasks))
		}
	})

	t.Run("create 201", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"text":"buy milk"}`))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
		}

		var task todo.Task
		if err := json.NewDecoder(rec.Body).Decode(&task); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if task.ID != 1 || task.Text != "buy milk" || task.Done {
			t.Fatalf("task = %+v, want id=1 text=buy milk done=false", task)
		}
	})

	t.Run("get found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}

		var task todo.Task
		if err := json.NewDecoder(rec.Body).Decode(&task); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if task.ID != 1 || task.Text != "buy milk" {
			t.Fatalf("task = %+v, want id=1 text=buy milk", task)
		}
	})

	t.Run("get 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tasks/999", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
		}
	})

	t.Run("done", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/tasks/1/done", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}

		var task todo.Task
		if err := json.NewDecoder(rec.Body).Decode(&task); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if !task.Done {
			t.Fatal("Done = false, want true")
		}
	})

	t.Run("delete 204", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
		}
	})

	t.Run("bad JSON 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{`))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})
}
