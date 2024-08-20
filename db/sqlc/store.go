package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	CreateAccountTx(ctx context.Context, arg AccountTxParams) (bool, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) error
}

type SQLStore struct {
	*Queries
	connPool *pgxpool.Pool
}

func NewStore(conn *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: conn,
		Queries:  New(conn),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.connPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rErr := tx.Rollback(ctx); rErr != nil {
			return fmt.Errorf("error tx::: %v", rErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
