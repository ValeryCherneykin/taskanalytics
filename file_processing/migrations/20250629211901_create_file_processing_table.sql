-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS uploaded_files (
    file_id TEXT PRIMARY KEY,
    file_name TEXT NOT NULL,
    file_path TEXT NOT NULL UNIQUE,
    size BIGINT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('pending', 'completed', 'failed')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS file_records (
    id SERIAL PRIMARY KEY,
    file_id TEXT NOT NULL REFERENCES uploaded_files(file_id),
    data TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_uploaded_files_status ON uploaded_files(status);
CREATE INDEX IF NOT EXISTS idx_file_records_file_id ON file_records(file_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_file_records_file_id;
DROP INDEX IF EXISTS idx_uploaded_files_status;
DROP TABLE IF EXISTS file_records;
DROP TABLE IF EXISTS uploaded_files;
-- +goose StatementEnd
