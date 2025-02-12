package store

import (
	"context"
	"time"
)

type Task struct {
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    int32     `json:"priority"`
	Id          int32     `json:"id"`
}

type TaskStore interface {
	CreateTask(ctx context.Context, title, description string, priority int32) (Task, error)
	GetTaskById(ctx context.Context, id int32) (Task, error)
	ListTasks(ctx context.Context) ([]Task, error)
	UpdateTask(ctx context.Context, id int32, title, description string, priority int32) (Task, error)
	DeleteTask(ctx context.Context, id int32) error
}
