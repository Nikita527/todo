package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type fileData struct {
	Next  int    `json:"next"`
	Tasks []Task `json:"tasks"`
}

func (s *Store) Save(path string) error {
	s.mu.Lock()
	tasks := make([]Task, len(s.tasks))
	copy(tasks, s.tasks)
	data := fileData{
		Next:  s.next,
		Tasks: tasks,
	}
	s.mu.Unlock()
	file, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	if err := os.WriteFile(path, file, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

func (s *Store) Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			s.mu.Lock()
			s.tasks = []Task{}
			s.next = 1
			s.mu.Unlock()
			return nil
		}
		return fmt.Errorf("failed to read file: %w", err)
	}
	var fileData fileData
	if err := json.Unmarshal(data, &fileData); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.next = fileData.Next
	s.tasks = fileData.Tasks
	return nil
}
