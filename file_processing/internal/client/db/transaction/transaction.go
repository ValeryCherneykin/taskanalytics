package transaction

import (
	"context"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/client/db"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/client/db/pg"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/logger"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type manager struct {
	db db.Transactor
}

func NewTransactionManager(db db.Transactor) db.TxManager {
	return &manager{
		db: db,
	}
}

func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, fn db.Handler) (err error) {
	logger.Debug("starting transaction", zap.Any("tx_options", opts))

	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if ok {
		logger.Debug("reusing existing transaction")
		return fn(ctx)
	}

	tx, err = m.db.BeginTx(ctx, opts)
	if err != nil {
		logger.Error("failed to begin transaction", zap.Error(err))
		return errors.Wrap(err, "can't begin transaction")
	}

	ctx = pg.MakeContextTx(ctx, tx)

	defer func() {
		if r := recover(); r != nil {
			logger.Error("panic recovered during transaction", zap.Any("recover", r))
			err = errors.Errorf("panic recovered: %v", r)
		}

		if err != nil {
			logger.Warn("rolling back transaction", zap.Error(err))
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				logger.Error("rollback failed", zap.Error(errRollback))
				err = errors.Wrapf(err, "errRollback: %v", errRollback)
			}
			return
		}

		logger.Debug("committing transaction")
		if commitErr := tx.Commit(ctx); commitErr != nil {
			logger.Error("failed to commit transaction", zap.Error(commitErr))
			err = errors.Wrap(commitErr, "tx commit failed")
		}
	}()

	if err = fn(ctx); err != nil {
		logger.Error("transaction failed inside function", zap.Error(err))
		err = errors.Wrap(err, "failed executing code inside transaction")
	}

	return err
}

func (m *manager) ReadCommitted(ctx context.Context, f db.Handler) error {
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.transaction(ctx, txOpts, f)
}
