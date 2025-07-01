package model

import "time"

type UploadedFile struct {
	FileID    int64
	FileName  string
	FilePath  string
	Size      int64
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type FileRecord struct {
	ID        int64
	FileID    string
	Data      string
	CreatedAt time.Time
}
