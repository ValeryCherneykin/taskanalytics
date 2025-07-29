package taskqueue

import (
	"context"

	"github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/model"
)

func (s *serv) GetTaskByID(ctx context.Context, taskID string) (*model.Task, error) {
	return s.queue.GetTaskByID(ctx, taskID)
}
