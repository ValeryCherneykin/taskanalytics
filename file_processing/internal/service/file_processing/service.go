package fileprocessing

import (
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/client/db"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/repository"
)

type serv struct {
	fileProcessingRepo repository.UploadedFileRepository
	txManager          db.TxManager
}

func NewService(fileProcessingRepo repository.UploadedFileRepository, txManager db.TxManager) *serv {
	return &serv{
		fileProcessingRepo: fileProcessingRepo,
		txManager:          txManager,
	}
}
