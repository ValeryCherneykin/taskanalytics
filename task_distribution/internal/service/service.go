package service

import (
	"context"
	"time"

	"github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/model"
)

type QueueService interface {
	AddTask(ctx context.Context, task *model.Task) error
	NextTask(ctx context.Context, timeout time.Duration) (*model.Task, error)
	GetTaskByID(ctx context.Context, taskID string) (*model.Task, error)
}
