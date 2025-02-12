package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/lohanguedes/taskify/internal/store"
)

func (api *Application) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Priority    int    `json:"priority"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	task := store.Task{
		Id:          1,
		Title:       payload.Title,
		Description: payload.Description,
		Priority:    int32(payload.Priority),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}
