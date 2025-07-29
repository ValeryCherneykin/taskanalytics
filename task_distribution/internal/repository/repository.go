package repository

import (
	"context"
	"time"

	"github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/model"
)

type TaskQueue interface {
	Enqueue(ctx context.Context, task *model.Task) error
	Dequeue(ctx context.Context, timeout time.Duration) (*model.Task, error)
	GetTaskByID(ctx context.Context, taskID string) (*model.Task, error)
}
