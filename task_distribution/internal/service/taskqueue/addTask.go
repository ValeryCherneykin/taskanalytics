package taskqueue

import (
	"context"
	"encoding/json"

	"github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/model"
)

func (s *serv) AddTask(ctx context.Context, task *model.Task) error {
	err := s.queue.Enqueue(ctx, task)
	if err != nil {
		return err
	}

	data, err := json.Marshal(task)
	if err != nil {
		return err
	}

	err = s.producer.SendMessage(ctx, "file-processing-tasks", []byte(task.TaskID), data)
	if err != nil {
		return err
	}

	return nil
}
