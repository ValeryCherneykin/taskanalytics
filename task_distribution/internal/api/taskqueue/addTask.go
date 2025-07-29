package taskqueue

import (
	"context"
	"time"

	"github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/model"
	desc "github.com/ValeryCherneykin/taskanalytics/task_distribution/pkg/task_distribution_v1"
	"github.com/ValeryCherneykin/taskanalytics/task_distribution/pkg/utils"
)

func (i *Implementation) SubmitFileProcessingTask(
	ctx context.Context,
	req *desc.SubmitFileProcessingTaskRequest,
) (*desc.SubmitFileProcessingTaskResponse, error) {
	task := &model.Task{
		TaskID:      utils.GenerateUUID(),
		FileID:      req.FileId,
		FilePath:    req.FilePath,
		FileName:    req.FileName,
		RecordCount: req.RecordCount,
		Status:      model.TaskStatusQueued,
		UploadedAt:  req.UploadedAt.AsTime(),
		UpdatedAt:   time.Now(),
		Message:     "",
	}

	err := i.taskQueueService.AddTask(ctx, task)
	if err != nil {
		return nil, err
	}

	return &desc.SubmitFileProcessingTaskResponse{
		TaskId: task.TaskID,
		Status: string(task.Status),
	}, nil
}
