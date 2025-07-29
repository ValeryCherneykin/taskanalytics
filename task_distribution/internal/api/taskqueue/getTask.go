package taskqueue

import (
	"context"

	"github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/converter"
	desc "github.com/ValeryCherneykin/taskanalytics/task_distribution/pkg/task_distribution_v1"
)

func (i *Implementation) GetTaskStatus(
	ctx context.Context,
	req *desc.GetTaskStatusRequest,
) (*desc.GetTaskStatusResponse, error) {
	task, err := i.taskQueueService.GetTaskByID(ctx, req.TaskId)
	if err != nil {
		return nil, err
	}

	return converter.ToProtoTask(task), nil
}
