package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeout  = 5 * time.Second
)

type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Builder squirrel.StatementBuilderType
	Pool    *pgxpool.Pool
}

func New(url string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(pg)
	}

	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = int32(pg.maxPoolSize)

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	ctx, cancel := context.WithTimeout(context.Background(), pg.connTimeout)
	defer cancel()

	if err := pg.Pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - ping failed: %w", err)
	}

	log.Println("Postgres successfully connected and pinged")

	return pg, nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

func (p *Postgres) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	if tx := extractTx(ctx); tx != nil {
		return tx.Exec(ctx, sql, args...)
	}
	return p.Pool.Exec(ctx, sql, args...)
}

func (p *Postgres) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if tx := extractTx(ctx); tx != nil {
		return tx.Query(ctx, sql, args...)
	}
	return p.Pool.Query(ctx, sql, args...)
}

func (p *Postgres) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if tx := extractTx(ctx); tx != nil {
		return tx.QueryRow(ctx, sql, args...)
	}
	return p.Pool.QueryRow(ctx, sql, args...)
}
