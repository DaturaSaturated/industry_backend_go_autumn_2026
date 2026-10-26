package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Task struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	UpdatedAt time.Time `json:"updatedAt"`
}
type TaskRepo interface {
	Create(title string) (Task, error)
	Get(id string) (Task, bool)
	List(done bool) []Task
	SetDone(id string, done bool) (Task, error)
}
type Clock interface{ Now() time.Time }

var (
	ErrNotFound     = errors.New("task not found")
	ErrInvalidTitle = errors.New("invalid title")
)

type inMemoryTaskRepo struct {
	mu    sync.RWMutex
	clock Clock
	seq   uint64
	tasks map[string]Task
}

func NewInMemoryTaskRepo(clock Clock) TaskRepo {
	tasks := make(map[string]Task)
	return &inMemoryTaskRepo{tasks: tasks, clock: clock}

}

func (r *inMemoryTaskRepo) Create(title string) (Task, error) {
	r.mu.Lock()
	var task Task
	task.ID = strconv.FormatUint(r.seq, 10)
	if strings.TrimSpace(title) == "" {
		r.mu.Unlock()
		return Task{}, ErrInvalidTitle
	}
	r.seq++
	task.Title = strings.TrimSpace(title)
	task.UpdatedAt = r.clock.Now()
	r.tasks[task.ID] = task
	r.mu.Unlock()
	return task, nil
}

func (r *inMemoryTaskRepo) Get(id string) (Task, bool) {
	r.mu.RLock() //на чтение ттолько, на запись не нужно
	defer r.mu.RUnlock()
	task, ok := r.tasks[id]
	if ok {
		return task, true
	}
	return Task{}, false
}

func (r *inMemoryTaskRepo) List(done bool) []Task {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]Task, 0)
	for _, t := range r.tasks {
		if t.Done == done {
			result = append(result, t)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].UpdatedAt.Equal(result[j].UpdatedAt) {
			return result[i].ID < result[j].ID
		}
		return result[i].UpdatedAt.After(result[j].UpdatedAt)
	})
	return result
}

func (r *inMemoryTaskRepo) SetDone(id string, done bool) (Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	task, ok := r.tasks[id]
	if !ok {
		return Task{}, ErrNotFound
	}

	if task.Done == done {
		return task, nil
	}
	task.Done = done
	task.UpdatedAt = r.clock.Now()
	r.tasks[id] = task
	return task, nil
}

type httpHandler struct{ repo TaskRepo }

func NewHTTPHandler(repo TaskRepo) http.Handler {
	return &httpHandler{repo: repo}
}
func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	path := r.URL.Path
	if method == "POST" {
		if path == "/tasks" {
			h.handleCreate(w, r)
			return
		}
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	if method == "GET" {
		if strings.HasPrefix(path, "/tasks/") {
			id := strings.TrimPrefix(r.URL.Path, "/tasks/")
			if id == "" || strings.Contains(id, "/") {
				writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
				return
			}
			h.handleGet(w, r, id)
			return
		}
		if path == "/tasks" {
			h.handleList(w, r)
			return
		}
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	if method == "PATCH" {
		if strings.HasPrefix(path, "/tasks/") {
			id := strings.TrimPrefix(r.URL.Path, "/tasks/")
			h.handlePatch(w, r, id)
			return
		}
		return
	}
	if path == "/tasks" || strings.HasPrefix(path, "/tasks/") {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
}

func (h *httpHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title *string `json:"title"`
	}

	if err := decodeStrictJSON(r.Body, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	if body.Title == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "title required"})
		return
	}

	task, err := h.repo.Create(*body.Title)
	if err != nil {
		if errors.Is(err, ErrInvalidTitle) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid title"})
			return
		}

		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusCreated, task)
}

func (h *httpHandler) handleGet(w http.ResponseWriter, r *http.Request, id string) {
	task, ok := h.repo.Get(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *httpHandler) handleList(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	values, exists := query["done"]
	done := false

	if exists {
		if len(values) != 1 || (values[0] != "true" && values[0] != "false") {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid done"})
			return
		}
		done = values[0] == "true"
	}
	tasks := h.repo.List(done)
	writeJSON(w, http.StatusOK, tasks)
}

func (h *httpHandler) handlePatch(w http.ResponseWriter, r *http.Request, id string) {
	var body struct {
		Done *bool `json:"done"`
	}
	if err := decodeStrictJSON(r.Body, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if body.Done == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "done required"})
		return
	}

	task, err := h.repo.SetDone(id, *body.Done)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func decodeStrictJSON(r io.Reader, v any) error {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	if err := dec.Decode(v); err != nil {
		return err
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("unexpected data after JSON")
	}

	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.Encode(v)
}
