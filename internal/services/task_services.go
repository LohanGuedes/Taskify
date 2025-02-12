package services

import (
	"context"

	"github.com/lohanguedes/taskify/internal/store"
)

type TaskService struct {
	Store store.TaskStore
}

func NewTaskService(store store.TaskStore) *TaskService {
	return &TaskService{Store: store}
}

func (s *TaskService) CreateTask(ctx context.Context, title, description string, priority int32) (store.Task, error) {
	// Add business logic here
	task, err := s.Store.CreateTask(ctx, title, description, priority)
	if err != nil {
		return store.Task{}, err
	}
	return task, err
}

func (s *TaskService) GetTask(ctx context.Context, id int32) (store.Task, error) {
	task, err := s.Store.GetTaskById(ctx, id)
	if err != nil {
		return store.Task{}, err
	}

	return task, err
}

func (s *TaskService) UpdateTask(ctx context.Context, id int32, title, description string, priority int32) (store.Task, error) {
	task, err := s.Store.UpdateTask(ctx, id, title, description, priority)
	if err != nil {
		return store.Task{}, err
	}

	return task, err
}

func (s *TaskService) DeleteTask(ctx context.Context, id int32) error {
	// Add business logic here...
	return s.Store.DeleteTask(ctx, id)
}
