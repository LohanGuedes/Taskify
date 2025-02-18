package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lohanguedes/taskify/internal/services"
	"github.com/lohanguedes/taskify/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestHandleCreateTask(t *testing.T) {
	api := Application{}

	payload := map[string]any{
		"title":       "Learn TDD",
		"description": "Get hands-on exp with TDD in Go!",
		"priority":    8000,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal("Failed to parse our request payload")
	}

	req := httptest.NewRequest("POST", "/api/v1/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(api.handleCreateTask)

	handler.ServeHTTP(rec, req)

	t.Logf("Rec body %s\n", rec.Body.Bytes())

	if rec.Code != http.StatusCreated {
		t.Errorf("Statuscode differs; got %d | want %d", rec.Code, http.StatusCreated)
	}

	var resBody map[string]any
	err = json.Unmarshal(rec.Body.Bytes(), &resBody)
	if err != nil {
		t.Fatalf("Failed to parse response body: %s\n", err.Error())
	}

	if resBody["title"] != payload["title"] {
		t.Errorf("title differs; got: %q | want: %q", resBody["title"], payload["title"])
	}
}

type MockTaskStore struct{}

func (mocktaskstore *MockTaskStore) CreateTask(ctx context.Context, title string, description string, priority int32) (store.Task, error) {
	return store.Task{
		Id:          1,
		Title:       title,
		Description: description,
		Priority:    priority,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func (mocktaskstore *MockTaskStore) GetTaskById(ctx context.Context, id int32) (store.Task, error) {
	return store.Task{
		Id:          id,
		Title:       "Mock Test Task",
		Description: "Mock Test Description",
		Priority:    1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func (mocktaskstore *MockTaskStore) ListTasks(ctx context.Context) ([]store.Task, error) {
	return []store.Task{
		{
			Id:          1,
			Title:       "Mock Test Task",
			Description: "Mock Test Description",
			Priority:    1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Id:          2,
			Title:       "Mock Test Task2",
			Description: "Mock Test Description2",
			Priority:    1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}, nil
}

func (mocktaskstore *MockTaskStore) UpdateTask(ctx context.Context, id int32, title string, description string, priority int32) (store.Task, error) {
	return store.Task{
		Id:          id,
		Title:       title,
		Description: description,
		Priority:    priority,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func (mocktaskstore *MockTaskStore) DeleteTask(ctx context.Context, id int32) error {
	return nil
}

func TestHandlerCreateTaskIntegration(t *testing.T) {
	mockStore := MockTaskStore{}
	taskService := services.NewTaskService(&mockStore)

	app := Application{
		TaskService: *taskService,
		Router:      chi.NewRouter(),
	}

	ts := httptest.NewServer(app.BindRoutes())
	defer ts.Close()

	payload := store.Task{
		Title:       "Integration Test Task",
		Description: "Testing the full API integration",
		Priority:    1,
	}

	body, err := json.Marshal(payload)
	assert.NoError(t, err, "Failed to Marshall request payload")

	resp, err := http.Post(ts.URL+"/api/v1/tasks", "application/json", bytes.NewReader(body))
	assert.NoError(t, err, "Failed  to send POST request")
	defer resp.Body.Close()

	assert.Equal(
		t,
		http.StatusCreated,
		resp.StatusCode,
		fmt.Sprintf("Expected status 201 got %d", resp.StatusCode),
	)

	var respTask store.Task

	err = json.NewDecoder(resp.Body).Decode(&respTask)
	assert.NoError(t, err, "Failed to decode response body")

	execptTask := store.Task{
		Id:          1,
		Title:       payload.Title,
		Description: payload.Description,
		Priority:    payload.Priority,
	}

	// Zero out cretion and update time for comparison
	respTask.CreatedAt = time.Time{}
	respTask.UpdatedAt = time.Time{}

	assert.Equal(t, execptTask, respTask, "The payload and response task do not match")
}

func TestHandleListTasks(t *testing.T) {
	mockStore := MockTaskStore{}
	taskService := services.NewTaskService(&mockStore)

	app := Application{
		TaskService: *taskService,
	}

	req := httptest.NewRequest("GET", "/api/v1/tasks", nil)
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(app.handleListTasks)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "Expected status 200 got %d", rec.Code)

	var tasks []store.Task
	err := json.NewDecoder(rec.Body).Decode(&tasks)
	assert.NoError(t, err, "Failed to decode response body")

	assert.Len(t, tasks, 2, "Expected 2 task got %d", len(tasks))
}

func TestHandleListTasksIntegration(t *testing.T) {
	mockStore := MockTaskStore{}
	taskService := services.NewTaskService(&mockStore)
	app := Application{
		TaskService: *taskService,
		Router:      chi.NewRouter(),
	}

	ts := httptest.NewServer(app.BindRoutes())
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/api/v1/tasks")
	assert.NoError(t, err, "Failed to send GET request")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status 200 got %d", resp.StatusCode)

	var tasks []store.Task
	err = json.NewDecoder(resp.Body).Decode(&tasks)
	assert.NoError(t, err, "Failed to decode response body")

	assert.Len(t, tasks, 2, "Expected 2 task got %d", len(tasks))

	expectedTask := store.Task{
		Id:          1,
		Title:       "Mock Test Task",
		Description: "Mock Test Description",
		Priority:    1,
	}
	tasks[0].CreatedAt = time.Time{}
	tasks[0].UpdatedAt = time.Time{}

	assert.Equal(t, expectedTask, tasks[0], "The payload and response task do not match")
}

func TestHandleGetTask(t *testing.T) {
	mockStore := MockTaskStore{}
	taskService := services.NewTaskService(&mockStore)
	app := Application{
		TaskService: *taskService,
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")

	req := httptest.NewRequest("GET", "/api/v1/tasks/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(app.handleGetTask)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "Expected status 200 got %d", rec.Code)

	var task store.Task
	err := json.NewDecoder(rec.Body).Decode(&task)
	assert.NoError(t, err, "Failed to decode response body")
	assert.Equal(t, int32(1), task.Id, "Expected task id 1 got %d", task.Id)
}

func TestHandleUpdateTask(t *testing.T) {
	mockStore := MockTaskStore{}
	taskService := services.NewTaskService(&mockStore)
	app := Application{
		TaskService: *taskService,
	}

	payload := map[string]any{
		"title":       "Updated Task",
		"description": "Updated Description",
		"priority":    1337,
	}

	body, err := json.Marshal(payload)
	assert.NoError(t, err, "Failed to Marshall request payload")
	t.Log(string(body))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")

	req := httptest.NewRequest("PUT", "/api/v1/tasks/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(app.handleUpdateTask)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "Expected status 200 got %d", rec.Code)

	var updatedTask store.Task
	err = json.NewDecoder(rec.Body).Decode(&updatedTask)
	assert.NoError(t, err, "Failed to decode response body")
	assert.Equal(t, "Updated Task", updatedTask.Title, "Expected title 'Updated Task' got %s", updatedTask.Title)
	assert.Equal(t, "Updated Description", updatedTask.Description, "Expected title 'Updated Task' got %s", updatedTask.Title)
	assert.Equal(t, int32(1337), updatedTask.Priority, "Expected title 'Updated Task' got %s", updatedTask.Title)
}

func TestHandleDeleteTask(t *testing.T) {
	mockStore := MockTaskStore{}
	taskService := services.NewTaskService(&mockStore)
	app := Application{
		TaskService: *taskService,
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")

	req := httptest.NewRequest("DELETE", "/api/v1/tasks/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(app.handleDeleteTask)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code, "Expected status 204 got %d", rec.Code)
}
