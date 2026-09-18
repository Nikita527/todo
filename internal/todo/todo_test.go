package todo_test

import (
	"errors"
	"testing"

	"github.com/Nikita527/todo/internal/todo"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		wantErr error
		wantID  int
	}{
		{
			name:   "success",
			text:   "buy milk",
			wantID: 1,
		},
		{
			name:    "empty text",
			text:    "",
			wantErr: todo.ErrEmptyText,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := todo.New()

			got, err := s.Add(tt.text)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Add(%q) error = %v, want %v", tt.text, err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			if got.ID != tt.wantID {
				t.Errorf("ID = %d, want %d", got.ID, tt.wantID)
			}
			if got.Text != tt.text {
				t.Errorf("Text = %q, want %q", got.Text, tt.text)
			}
			if got.Done {
				t.Error("Done = true, want false")
			}
		})
	}
}

func TestListIsCopy(t *testing.T) {
	s := todo.New()
	if _, err := s.Add("original"); err != nil {
		t.Fatalf("Add: %v", err)
	}

	list := s.List()
	if len(list) != 1 {
		t.Fatalf("List len = %d, want 1", len(list))
	}

	list[0].Text = "mutated"
	list[0].Done = true

	again := s.List()
	if again[0].Text != "original" {
		t.Errorf("store Text = %q, want %q (List must return a copy)", again[0].Text, "original")
	}
	if again[0].Done {
		t.Error("store Done = true, want false (List must return a copy)")
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*todo.Store) int
		id      int
		wantErr error
		want    todo.Task
	}{
		{
			name: "found",
			setup: func(s *todo.Store) int {
				task, err := s.Add("found task")
				if err != nil {
					t.Fatalf("Add: %v", err)
				}
				return task.ID
			},
			want: todo.Task{ID: 1, Text: "found task", Done: false},
		},
		{
			name: "not found",
			setup: func(s *todo.Store) int {
				if _, err := s.Add("other"); err != nil {
					t.Fatalf("Add: %v", err)
				}
				return 999
			},
			id:      999,
			wantErr: todo.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := todo.New()
			id := tt.id
			if tt.setup != nil {
				id = tt.setup(s)
			}

			got, err := s.Get(id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Get(%d) error = %v, want %v", id, err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			if got != tt.want {
				t.Errorf("Get(%d) = %+v, want %+v", id, got, tt.want)
			}
		})
	}
}

func TestDone(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*todo.Store) int
		wantErr error
	}{
		{
			name: "success",
			setup: func(s *todo.Store) int {
				task, err := s.Add("mark done")
				if err != nil {
					t.Fatalf("Add: %v", err)
				}
				return task.ID
			},
		},
		{
			name: "not found",
			setup: func(*todo.Store) int {
				return 42
			},
			wantErr: todo.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := todo.New()
			id := tt.setup(s)

			err := s.Done(id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Done(%d) error = %v, want %v", id, err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}

			got, err := s.Get(id)
			if err != nil {
				t.Fatalf("Get after Done: %v", err)
			}
			if !got.Done {
				t.Error("Done = false after Done(), want true")
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*todo.Store) int
		wantErr error
	}{
		{
			name: "success",
			setup: func(s *todo.Store) int {
				task, err := s.Add("to delete")
				if err != nil {
					t.Fatalf("Add: %v", err)
				}
				return task.ID
			},
		},
		{
			name: "not found",
			setup: func(*todo.Store) int {
				return 42
			},
			wantErr: todo.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := todo.New()
			id := tt.setup(s)

			err := s.Delete(id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Delete(%d) error = %v, want %v", id, err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}

			if _, err := s.Get(id); !errors.Is(err, todo.ErrNotFound) {
				t.Fatalf("Get after Delete: error = %v, want %v", err, todo.ErrNotFound)
			}
		})
	}
}
