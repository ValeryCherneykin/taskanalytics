package taskqueue

import (
	"github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/service"
	desc "github.com/ValeryCherneykin/taskanalytics/task_distribution/pkg/task_distribution_v1"
)

type Implementation struct {
	desc.UnimplementedTaskServiceServer
	taskQueueService service.QueueService
}

func NewImplementation(taskQueueService service.QueueService) *Implementation {
	return &Implementation{
		taskQueueService: taskQueueService,
	}
}
