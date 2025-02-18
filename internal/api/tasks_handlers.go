package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
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

func (api *Application) handleListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := api.TaskService.Store.ListTasks(r.Context())
	if err != nil {
		http.Error(w, "Failed to list tasks, try again later", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (api *Application) handleGetTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Id must be a valid integer", http.StatusBadRequest)
		return
	}

	task, err := api.TaskService.Store.GetTaskById(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "task Not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (api *Application) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Id must be a valid integer", http.StatusBadRequest)
		return
	}

	var payload store.Task
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	task, err := api.TaskService.UpdateTask(r.Context(), int32(id), payload.Title, payload.Description, payload.Priority)
	if err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (api *Application) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Id must be a valid integer", http.StatusBadRequest)
		return
	}

	err = api.TaskService.Store.DeleteTask(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
