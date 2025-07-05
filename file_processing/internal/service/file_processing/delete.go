package fileprocessing

import (
	"context"
	"fmt"
)

func (s *serv) Delete(ctx context.Context, id int64) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		errTx := s.fileProcessingRepository.Delete(ctx, id)
		if errTx != nil {
			return errTx
		}

		_, err := s.fileProcessingRepository.Get(ctx, id)
		if err == nil {
			return fmt.Errorf("file with id %d still exists after delete", id)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

