package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"example.com/m/v2/internal/database"
	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	store *database.TaskStore
}

func NewTaskHandler(store *database.TaskStore) *TaskHandler {
	return &TaskHandler{store: store}
}

func (h *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.GetAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get tasks")
		return
	}

	writeJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) Sleep(w http.ResponseWriter, r *http.Request) {
	log.Println("sleep started")
	time.Sleep(8 * time.Second)
	log.Println("sleep finished")
	_, _ = w.Write([]byte("ok"))
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "failed to encode", http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func (h *TaskHandler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", h.GetAll)
	r.Get("/sleep", h.Sleep)
	return r
}
