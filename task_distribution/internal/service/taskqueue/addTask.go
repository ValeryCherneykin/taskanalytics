package taskqueue

import (
	"context"

	"github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/model"
)

func (s *serv) AddTask(ctx context.Context, task *model.Task) error {
	return s.queue.Enqueue(ctx, task)
}
