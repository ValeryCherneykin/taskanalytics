package taskqueue

import (
	"context"
	"time"

	"github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/model"
)

func (s *serv) NextTask(ctx context.Context, timeout time.Duration) (*model.Task, error) {
	return s.queue.Dequeue(ctx, timeout)
}
