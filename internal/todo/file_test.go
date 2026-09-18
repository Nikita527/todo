package todo_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Nikita527/todo/internal/todo"
)

func TestSaveLoad(t *testing.T) {
	t.Run("round-trip", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "todos.json")

		s := todo.New()
		t1, err := s.Add("buy milk")
		if err != nil {
			t.Fatalf("Add: %v", err)
		}
		t2, err := s.Add("walk dog")
		if err != nil {
			t.Fatalf("Add: %v", err)
		}
		if err := s.Done(t1.ID); err != nil {
			t.Fatalf("Done: %v", err)
		}

		if err := s.Save(path); err != nil {
			t.Fatalf("Save: %v", err)
		}

		loaded := todo.New()
		if err := loaded.Load(path); err != nil {
			t.Fatalf("Load: %v", err)
		}

		list := loaded.List()
		if len(list) != 2 {
			t.Fatalf("List len = %d, want 2", len(list))
		}
		if list[0].ID != t1.ID || list[0].Text != "buy milk" || !list[0].Done {
			t.Errorf("List[0] = %+v, want ID=%d Text=%q Done=true", list[0], t1.ID, "buy milk")
		}
		if list[1].ID != t2.ID || list[1].Text != "walk dog" || list[1].Done {
			t.Errorf("List[1] = %+v, want ID=%d Text=%q Done=false", list[1], t2.ID, "walk dog")
		}

		next, err := loaded.Add("new task")
		if err != nil {
			t.Fatalf("Add after Load: %v", err)
		}
		if next.ID != 3 {
			t.Errorf("Add after Load ID = %d, want 3", next.ID)
		}
	})

	t.Run("missing file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "no-such.json")
		s := todo.New()
		if _, err := s.Add("should disappear"); err != nil {
			t.Fatalf("Add: %v", err)
		}

		if err := s.Load(path); err != nil {
			t.Fatalf("Load missing file: %v, want nil", err)
		}
		if len(s.List()) != 0 {
			t.Errorf("List = %v, want empty", s.List())
		}
	})

	t.Run("bad json", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "bad.json")
		if err := os.WriteFile(path, []byte("{not-json"), 0o644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}

		s := todo.New()
		if err := s.Load(path); err == nil {
			t.Fatal("Load bad json: err = nil, want error")
		}
	})
}
