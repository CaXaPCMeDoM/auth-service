package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type txKey struct{}

func (p *Postgres) WithinTransaction(ctx context.Context, f func(ctx context.Context) error) error {
	tx, err := p.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	ctxWithTx := injectTx(ctx, &tx)

	if err := f(ctxWithTx); err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("tx func error: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit error: %w", err)
	}

	return nil
}

func injectTx(ctx context.Context, tx *pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func extractTx(ctx context.Context) pgx.Tx {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}
	return nil
}
