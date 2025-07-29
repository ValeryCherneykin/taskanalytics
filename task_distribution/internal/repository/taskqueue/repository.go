package taskqueue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/client/taskstate"
	"github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/model"
	"github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/repository"
)

const queueName = "task_queue"

type repo struct {
	redis taskstate.RedisClient
}

func NewRepository(redisClient taskstate.RedisClient) repository.TaskQueue {
	return &repo{
		redis: redisClient,
	}
}

func (r *repo) Enqueue(ctx context.Context, task *model.Task) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}

	if err := r.redis.Set(ctx, "task:"+task.TaskID, data); err != nil {
		return err
	}

	return r.redis.LPush(ctx, queueName, data)
}

func (r *repo) Dequeue(ctx context.Context, timeout time.Duration) (*model.Task, error) {
	raw, err := r.redis.BRPop(ctx, queueName, timeout)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}

	bytes, ok := raw.([]byte)
	if !ok {
		return nil, ErrInvalidTaskFormat
	}

	var task model.Task
	if err := json.Unmarshal(bytes, &task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *repo) GetTaskByID(ctx context.Context, taskID string) (*model.Task, error) {
	raw, err := r.redis.Get(ctx, "task:"+taskID)
	if err != nil || raw == nil {
		return nil, err
	}

	bytes, ok := raw.([]byte)
	if !ok {
		return nil, ErrInvalidTaskFormat
	}

	var task model.Task
	if err := json.Unmarshal(bytes, &task); err != nil {
		return nil, err
	}

	return &task, nil
}
