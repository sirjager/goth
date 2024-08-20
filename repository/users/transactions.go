package users

import (
	"context"
	"fmt"
)

func (r *userRepo) ExecTx(ctx context.Context, fn func() error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	if err = fn(); err != nil {
		if rollBackErr := tx.Rollback(ctx); rollBackErr != nil {
			return fmt.Errorf("transaction error: %v , rollback error: %v", err, rollBackErr)
		}
		return err
	}
	return tx.Commit(ctx)
}
