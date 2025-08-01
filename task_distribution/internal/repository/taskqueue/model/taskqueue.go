package model

import "time"

type TaskStatus string

const (
	TaskStatusUnspecified TaskStatus = "UNSPECIFIED"
	TaskStatusQueued      TaskStatus = "QUEUED"
	TaskStatusInProgress  TaskStatus = "IN_PROGRESS"
	TaskStatusCompleted   TaskStatus = "COMPLETED"
	TaskStatusFailed      TaskStatus = "FAILED"
)

type Task struct {
	TaskID      string     `json:"task_id"`
	FileID      int64      `json:"file_id"`
	FilePath    string     `json:"file_path"`
	FileName    string     `json:"file_name"`
	RecordCount int64      `json:"record_count"`
	Status      TaskStatus `json:"status"`
	UploadedAt  time.Time  `json:"uploaded_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Message     string     `json:"message"`
}
