package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TXError struct {
	Err         error
	RollbackErr error
	CommitErr   error
}

func (e *TXError) Error() string {
	return fmt.Sprintf("err: %v; rollback error: %v; commit error: %v", e.Err, e.RollbackErr, e.CommitErr)
}

type DAO interface {
	Querier
	DBTX
	GetConn() *pgxpool.Pool
	ExecTx(ctx context.Context, fn func(DAO) error) *TXError
	QueryTx(ctx context.Context, fn func(DAO) (interface{}, error)) (interface{}, *TXError)
	Close()
}

// extend Queries struct ability
type dao struct {
	*Queries
	db *pgxpool.Pool
}

func NewDAO(dbUrl string) (DAO, error) {
	// 5 seconds timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return nil, err
	}

	return &dao{
		Queries: New(conn),
		db:      conn,
	}, nil
}

func (d *dao) GetConn() *pgxpool.Pool {
	return d.db
}

func (d *dao) Exec(ctx context.Context, query string, params ...interface{}) (pgconn.CommandTag, error) {
	return d.db.Exec(ctx, query, params...)
}

func (d *dao) Query(ctx context.Context, query string, params ...interface{}) (pgx.Rows, error) {
	return d.db.Query(ctx, query, params...)
}

func (d *dao) QueryRow(ctx context.Context, query string, params ...interface{}) pgx.Row {
	return d.db.QueryRow(ctx, query, params...)
}

func (d *dao) ExecTx(ctx context.Context, fn func(DAO) error) *TXError {
	tx, err := d.db.Begin(ctx)
	if err != nil {
		return &TXError{Err: err}
	}

	if err = fn(d); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return &TXError{Err: err, RollbackErr: rbErr}
		}
		return &TXError{Err: err}
	}

	cErr := tx.Commit(ctx)
	if cErr != nil {
		return &TXError{CommitErr: cErr}
	}

	return nil
}

func (d *dao) QueryTx(ctx context.Context, fn func(DAO) (interface{}, error)) (interface{}, *TXError) {
	tx, err := d.db.Begin(ctx)
	if err != nil {
		return nil, &TXError{Err: err}
	}

	res, err := fn(d)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return nil, &TXError{Err: err, RollbackErr: rbErr}
		}
		return nil, &TXError{Err: err}
	}

	cErr := tx.Commit(ctx)
	if cErr != nil {
		return nil, &TXError{CommitErr: cErr}
	}

	return res, nil
}

func (d *dao) Close() {
	d.db.Close()
}
