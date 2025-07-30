package taskqueue

import (
	"github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/client/kafka"
	"github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/repository"
	"github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/service"
)

type serv struct {
	queue    repository.TaskQueue
	producer kafka.Producer
}

func NewService(queue repository.TaskQueue, producer kafka.Producer) service.QueueService {
	return &serv{
		queue:    queue,
		producer: producer,
	}
}
