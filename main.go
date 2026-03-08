package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

type Task struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type TaskStore struct {
	sync.RWMutex
	tasks map[string]Task
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks: make(map[string]Task),
	}
}

func (s *TaskStore) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		s.RLock()
		list := make([]Task, 0, len(s.tasks))
		for _, v := range s.tasks {
			list = append(list, v)
		}
		s.RUnlock()
		json.NewEncoder(w).Encode(list)
	case http.MethodPost:
		var t Task
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil || t.ID == "" {
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}
		s.Lock()
		s.tasks[t.ID] = t
		s.Unlock()
		w.WriteHeader(http.StatusCreated)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	server := NewTaskStore()
	http.ListenAndServe(":8080", server)
}
