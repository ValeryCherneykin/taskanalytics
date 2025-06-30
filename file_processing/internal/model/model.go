package model

import (
	"time"

	"gorm.io/gorm"
)

type UploadedFile struct {
	FileID    string    `gorm:"type:text;primaryKey"`
	FileName  string    `gorm:"type:text;not null"`
	FilePath  string    `gorm:"type:text;not null;unique"`
	Size      int64     `gorm:"not null"`
	Status    string    `gorm:"type:text;not null;check:status IN ('pending','completed','failed')"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type FileRecord struct {
	ID        uint      `gorm:"primaryKey"`
	FileID    string    `gorm:"type:text;index;not null"`
	Data      string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"not null"`
}
