package model

import "time"

type UploadedFile struct {
	FileID    string
	FileName  string
	FilePath  string
	Size      int64
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type FileRecord struct {
	ID        uint
	FileID    string
	Data      string
	CreatedAt time.Time
}
