package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Nikita527/todo/internal/todo"
)

type API struct {
	store *todo.Store
	file  string
}

func New(store *todo.Store, file string) *API {
	return &API{store: store, file: file}
}

func (a *API) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", a.listTasks)
	mux.HandleFunc("POST /tasks", a.createTask)
	mux.HandleFunc("GET /tasks/{id}", a.getTask)
	mux.HandleFunc("POST /tasks/{id}/done", a.doneTask)
	mux.HandleFunc("DELETE /tasks/{id}", a.deleteTask)
	return mux
}

func (a *API) listTasks(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, a.store.List())
}

func (a *API) createTask(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	task, err := a.store.Add(body.Text)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if !a.persist(w) {
		return
	}
	writeJSON(w, http.StatusCreated, task)
}

func (a *API) getTask(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	task, err := a.store.Get(id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (a *API) doneTask(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	if err := a.store.Done(id); err != nil {
		writeStoreError(w, err)
		return
	}

	task, err := a.store.Get(id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if !a.persist(w) {
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (a *API) deleteTask(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	if err := a.store.Delete(id); err != nil {
		writeStoreError(w, err)
		return
	}
	if !a.persist(w) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func pathID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

func writeStoreError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, todo.ErrEmptyText):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, todo.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}

func (a *API) persist(w http.ResponseWriter) bool {
	if err := a.store.Save(a.file); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return false
	}
	return true
}
