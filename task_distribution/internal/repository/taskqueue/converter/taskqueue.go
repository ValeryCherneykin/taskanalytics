package converter

import (
	"github.com/ValeryCherneykin/taskanalytics/task_distribution/internal/repository/taskqueue/model"
	desc "github.com/ValeryCherneykin/taskanalytics/task_distribution/pkg/task_distribution_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToProtoTaskStatus(status model.TaskStatus) desc.TaskStatus {
	switch status {
	case model.TaskStatusQueued:
		return desc.TaskStatus_TASK_STATUS_QUEUED
	case model.TaskStatusInProgress:
		return desc.TaskStatus_TASK_STATUS_IN_PROGRESS
	case model.TaskStatusCompleted:
		return desc.TaskStatus_TASK_STATUS_COMPLETED
	case model.TaskStatusFailed:
		return desc.TaskStatus_TASK_STATUS_FAILED
	default:
		return desc.TaskStatus_TASK_STATUS_UNSPECIFIED
	}
}

func ToModelTaskStatus(status desc.TaskStatus) model.TaskStatus {
	switch status {
	case desc.TaskStatus_TASK_STATUS_QUEUED:
		return model.TaskStatusQueued
	case desc.TaskStatus_TASK_STATUS_IN_PROGRESS:
		return model.TaskStatusInProgress
	case desc.TaskStatus_TASK_STATUS_COMPLETED:
		return model.TaskStatusCompleted
	case desc.TaskStatus_TASK_STATUS_FAILED:
		return model.TaskStatusFailed
	default:
		return model.TaskStatusUnspecified
	}
}

func ToProtoTask(task *model.Task) *desc.GetTaskStatusResponse {
	return &desc.GetTaskStatusResponse{
		TaskId:    task.TaskID,
		Status:    ToProtoTaskStatus(task.Status),
		Message:   task.Message,
		UpdatedAt: timestamppb.New(task.UpdatedAt),
	}
}

func ToModelTask(resp *desc.GetTaskStatusResponse) *model.Task {
	return &model.Task{
		TaskID:    resp.TaskId,
		Status:    ToModelTaskStatus(resp.Status),
		Message:   resp.Message,
		UpdatedAt: resp.UpdatedAt.AsTime(),
	}
}
