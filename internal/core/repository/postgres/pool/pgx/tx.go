package core_pgx_pool

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type txKey struct{}

func txFromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}

func (p *Pool) WithinTx(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	if _, ok := txFromContext(ctx); ok {
		return fn(ctx)
	}

	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", mapErrors(err))
	}
	defer func() {
		_ = tx.Rollback(context.WithoutCancel(ctx))
	}()

	if err := fn(context.WithValue(ctx, txKey{}, tx)); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", mapErrors(err))
	}

	return nil
}
