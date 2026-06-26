package repository

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
)

type Store struct {
	db bun.IDB
}

func NewStore(db bun.IDB) *Store {
	if db == nil {
		panic("repository.NewStore: db is nil")
	}
	return &Store{db: db}
}

func (s *Store) RunInTx(ctx context.Context, fn func(context.Context, *Store) error) error {
	if runner, ok := s.db.(interface {
		RunInTx(context.Context, *sql.TxOptions, func(context.Context, bun.Tx) error) error
	}); ok {
		return runner.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
			return fn(ctx, NewStore(tx))
		})
	}
	return fn(ctx, s)
}
