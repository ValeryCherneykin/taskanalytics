package model

import "gorm.io/gorm"

type UploadedFile struct {
	gorm.Model
	FileName string `gorm:"not null"`
	FilePath string `gorm:"not null;unique"`
}
