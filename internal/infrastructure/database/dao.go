package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type TXError struct {
	Err   error
	RbErr error
	CErr  error
}

func (e *TXError) Error() string {
	return fmt.Sprintf("err: %v; rollback error: %v; commit error: %v", e.Err, e.RbErr, e.CErr)
}

type DAO interface {
	Querier
	DBTX
	GetConn() *sql.DB
	ExecTx(ctx context.Context, fn func(DAO) error) *TXError
	QueryTx(ctx context.Context, fn func(DAO) (interface{}, error)) (interface{}, *TXError)
	Close() error
}

// extend Queries struct ability
type dao struct {
	*Queries
	db *sql.DB
}

func NewDAO(dbUrl string) (DAO, error) {
	conn, err := sql.Open("pgx", dbUrl)
	if err != nil {
		return nil, err
	}

	// 5 seconds timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = conn.PingContext(ctx); err != nil {
		return nil, err
	}

	return &dao{
		Queries: New(conn),
		db:      conn,
	}, nil
}

func (d *dao) GetConn() *sql.DB {
	return d.db
}

func (d *dao) ExecContext(ctx context.Context, stmt string, args ...interface{}) (sql.Result, error) {
	return d.db.ExecContext(ctx, stmt, args...)
}

func (d *dao) PrepareContext(ctx context.Context, stmt string) (*sql.Stmt, error) {
	return d.db.PrepareContext(ctx, stmt)
}

func (d *dao) QueryContext(ctx context.Context, stmt string, args ...interface{}) (*sql.Rows, error) {
	return d.db.QueryContext(ctx, stmt, args...)
}

func (d *dao) QueryRowContext(ctx context.Context, stmt string, args ...interface{}) *sql.Row {
	return d.db.QueryRowContext(ctx, stmt, args...)
}

func (d *dao) ExecTx(ctx context.Context, fn func(DAO) error) *TXError {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return &TXError{Err: err}
	}

	if err = fn(d); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return &TXError{Err: err, RbErr: rbErr}
		}
		return &TXError{Err: err}
	}

	cErr := tx.Commit()
	if cErr != nil {
		return &TXError{CErr: cErr}
	}

	return nil
}

func (d *dao) QueryTx(ctx context.Context, fn func(DAO) (interface{}, error)) (interface{}, *TXError) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, &TXError{Err: err}
	}

	res, err := fn(d)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return nil, &TXError{Err: err, RbErr: rbErr}
		}
		return nil, &TXError{Err: err}
	}

	cErr := tx.Commit()
	if cErr != nil {
		return nil, &TXError{CErr: cErr}
	}

	return res, nil
}

func (d *dao) Close() error {
	return d.db.Close()
}
