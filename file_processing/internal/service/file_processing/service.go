package fileprocessing

import (
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/client/db"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/repository"
)

type serv struct {
	fileProcessingRepository repository.UploadedFileRepository
	txManager                db.TxManager
}

func NewService(
	fileProcessingRepository repository.UploadedFileRepository,
	txManager db.TxManager,
) *serv {
	return &serv{
		fileProcessingRepository: fileProcessingRepository,
		txManager:                txManager,
	}
}
