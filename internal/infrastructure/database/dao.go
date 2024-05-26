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
	ExecTx(ctx context.Context, txOptions any, fn func(DAO) error) *TXError
	QueryTx(ctx context.Context, txOptions any, fn func(DAO) (interface{}, error)) (interface{}, *TXError)
	QueryBatch(ctx context.Context, queries []BatchedQuery) error
	QueryRowBatch(ctx context.Context, queries []BatchedQueryRow) error
	ExecBatch(ctx context.Context, queries []BatchedExec) error
	Close()
}

// extend Queries struct ability
type postgresDAO struct {
	*Queries
	db   *pgxpool.Pool
	dbtx *pgx.Tx
}

func NewPostgresDAO(dbUrl string) (DAO, error) {
	// 5 seconds timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return nil, err
	}

	return &postgresDAO{
		Queries: New(conn),
		db:      conn,
	}, nil
}

func (d *postgresDAO) GetConn() *pgxpool.Pool {
	return d.db
}

func (d *postgresDAO) Exec(ctx context.Context, query string, params ...interface{}) (pgconn.CommandTag, error) {
	if d.dbtx != nil {
		return (*d.dbtx).Exec(ctx, query, params...)
	}
	return d.db.Exec(ctx, query, params...)
}

func (d *postgresDAO) Query(ctx context.Context, query string, params ...interface{}) (pgx.Rows, error) {
	if d.dbtx != nil {
		return (*d.dbtx).Query(ctx, query, params...)
	}
	return d.db.Query(ctx, query, params...)
}

func (d *postgresDAO) QueryRow(ctx context.Context, query string, params ...interface{}) pgx.Row {
	if d.dbtx != nil {
		return (*d.dbtx).QueryRow(ctx, query, params...)
	}
	return d.db.QueryRow(ctx, query, params...)
}

func (d *postgresDAO) ExecTx(ctx context.Context, txOptions any, fn func(DAO) error) *TXError {
	opts, ok := txOptions.(*pgx.TxOptions)
	var (
		tx  pgx.Tx
		err error
	)
	if !ok || opts == nil {
		tx, err = d.db.Begin(ctx)
	} else {
		tx, err = d.db.BeginTx(ctx, *opts)
	}

	if err != nil {
		return &TXError{Err: err}
	}

	q := &postgresDAO{
		Queries: New(tx),
		db:      d.db,
		dbtx:    &tx,
	}
	if err = fn(q); err != nil {
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

func (d *postgresDAO) QueryTx(ctx context.Context, txOptions any, fn func(DAO) (interface{}, error)) (interface{}, *TXError) {
	opts, ok := txOptions.(*pgx.TxOptions)
	var (
		tx  pgx.Tx
		err error
	)
	if !ok || opts == nil {
		tx, err = d.db.Begin(ctx)
	} else {
		tx, err = d.db.BeginTx(ctx, *opts)
	}

	if err != nil {
		return nil, &TXError{Err: err}
	}

	q := &postgresDAO{
		Queries: New(tx),
		db:      d.db,
		dbtx:    &tx,
	}
	res, err := fn(q)
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

// Batch Query

type BatchedQuery struct {
	SQL    string
	Params []interface{}
	Fn     func(row pgx.Rows) error
}

func (d *postgresDAO) QueryBatch(ctx context.Context, queries []BatchedQuery) error {
	batch := &pgx.Batch{}
	for _, q := range queries {
		batch.Queue(q.SQL, q.Params...).Query(q.Fn)
	}

	return d.db.SendBatch(ctx, batch).Close()
}

type BatchedQueryRow struct {
	SQL    string
	Params []interface{}
	Fn     func(row pgx.Row) error
}

func (d *postgresDAO) QueryRowBatch(ctx context.Context, queries []BatchedQueryRow) error {
	batch := &pgx.Batch{}
	for _, q := range queries {
		batch.Queue(q.SQL, q.Params...).QueryRow(q.Fn)
	}

	return d.db.SendBatch(ctx, batch).Close()
}

// Batch Exec

type BatchedExec struct {
	SQL    string
	Params []interface{}
	Fn     func(ct pgconn.CommandTag) error
}

func (d *postgresDAO) ExecBatch(ctx context.Context, queries []BatchedExec) error {
	batch := &pgx.Batch{}
	for _, q := range queries {
		batch.Queue(q.SQL, q.Params...).Exec(q.Fn)
	}

	return d.db.SendBatch(ctx, batch).Close()
}

func (d *postgresDAO) Close() {
	d.db.Close()
}
