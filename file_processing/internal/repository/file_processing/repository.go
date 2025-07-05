package fileprocessing

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/client/db"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/model"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/repository"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/repository/file_processing/converter"
	modelRepo "github.com/ValeryCherneykin/taskanalytics/file_processing/internal/repository/file_processing/model"
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

func NewRepository(db db.Client) repository.UploadedFileRepository {
	return &repo{
		db: db,
	}
}

func (r *repo) Create(ctx context.Context, file *model.UploadedFile) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(fileNameColumn, filePathColumn, sizeColumn, statusColumn).
		Values(file.FileName, file.FilePath, file.Size, file.Status).
		Suffix("RETURNING ID")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "file_processing.Create",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.UploadedFile, error) {
	builder := sq.Select(idColumn, fileNameColumn, filePathColumn, sizeColumn, statusColumn, createdAtColumn, updatedAtColumn, deletedAtColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "file_processing.Get",
		QueryRaw: query,
	}

	var file modelRepo.UploadedFile
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&file.FileID, &file.FileName, &file.FilePath, &file.Size, &file.Status, &file.CreatedAt, &file.UpdatedAt, &file.DeletedAt)

	return converter.ToFileMetadataFromRepo(&file), err
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	builder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "file_processing.Delete",
		QueryRaw: query,
	}
	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return err
}

func (r *repo) Update(ctx context.Context, file *model.UploadedFile) error {
	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(fileNameColumn, file.FileName).
		Set(filePathColumn, file.FilePath).
		Set(sizeColumn, file.Size).
		Set(statusColumn, file.Status).
		Where(sq.Eq{idColumn: file.FileID})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "file_processing.Update",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return err
}

func (r *repo) List(ctx context.Context, limit, offset uint64) ([]*model.UploadedFile, error) {
	builder := sq.Select(
		idColumn, fileNameColumn, filePathColumn, sizeColumn,
		statusColumn, createdAtColumn, updatedAtColumn, deletedAtColumn,
	).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		OrderBy(createdAtColumn + " DESC").
		Limit(limit).
		Offset(offset)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "file_processing.List",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*model.UploadedFile
	for rows.Next() {
		var file modelRepo.UploadedFile
		if err := rows.Scan(
			&file.FileID, &file.FileName, &file.FilePath,
			&file.Size, &file.Status, &file.CreatedAt,
			&file.UpdatedAt, &file.DeletedAt,
		); err != nil {
			return nil, err
		}
		files = append(files, converter.ToFileMetadataFromRepo(&file))
	}

	return files, nil
}
