package taskqueue

import (
	"github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/repository"
	"github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/service"
)

type serv struct {
	queue repository.TaskQueue
}

func NewService(queue repository.TaskQueue) service.QueueService {
	return &serv{
		queue: queue,
	}
}
