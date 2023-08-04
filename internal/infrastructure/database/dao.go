package database

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DAO interface {
	Querier
	GetConn() *sql.DB
	Close() error
}

// extend Queries struct ability
type dao struct {
	*Queries
	conn *sql.DB
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
		conn:    conn,
	}, nil
}

func (d *dao) GetConn() *sql.DB {
	return d.conn
}

func (d *dao) Close() error {
	return d.conn.Close()
}
