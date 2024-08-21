package db

import (
	"context"
	"fmt"
)

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return nil
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		rbErr := tx.Rollback(ctx)
		if rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}
