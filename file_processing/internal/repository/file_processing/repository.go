package fileprocessing

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/client/db"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/repository/file_processing/model"
)

const (
	tableName = "uploaded_files"

	idColumn        = "file_id"
	fileNameColumn  = "file_name"
	filePathColumn  = "file_path"
	sizeColumn      = "size"
	statusColumn    = "status"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
	deletedAtColumn = "deleted_at"
)

type repo struct {
	db db.Client
}

// func NewRepository(db db.Client) repository.UploadedFileRepository {
// 	return &repo{
// 		db: db,
// 	}
// }

func (r *repo) Create(ctx context.Context, info model.UploadedFile) (int64, error) {
	builder := squirrel.Insert(tableName).
		PlaceholderFormat(squirrel.Dollar).
		Columns(fileNameColumn, filePathColumn).
		Values(info.FileName, info.FilePath).
		Suffix("RETURNING file_id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "file_processing_repository.Create",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (model.UploadedFile, error)
