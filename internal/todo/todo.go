package todo

import (
	"errors"
	"sync"
)

var ErrEmptyText = errors.New("text is empty")
var ErrNotFound = errors.New("task not found")

type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type Store struct {
	mu    sync.Mutex
	tasks []Task
	next  int
}

func New() *Store {
	return &Store{
		tasks: make([]Task, 0),
		next:  1,
	}
}

func (s *Store) Add(text string) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if text == "" {
		return Task{}, ErrEmptyText
	}
	s.tasks = append(s.tasks, Task{
		ID:   s.next,
		Text: text,
		Done: false,
	})
	s.next++
	return s.tasks[len(s.tasks)-1], nil
}

func (s *Store) List() []Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	tasks := make([]Task, len(s.tasks))
	copy(tasks, s.tasks)
	return tasks
}

func (s *Store) Done(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, task := range s.tasks {
		if task.ID == id {
			s.tasks[i].Done = true
			return nil
		}
	}
	return ErrNotFound
}

func (s *Store) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, task := range s.tasks {
		if task.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}

func (s *Store) Get(id int) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, task := range s.tasks {
		if task.ID == id {
			return task, nil
		}
	}
	return Task{}, ErrNotFound
}
