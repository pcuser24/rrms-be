package database

import (
	"database/sql"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	ddsqlx "gopkg.in/DataDog/dd-trace-go.v1/contrib/jmoiron/sqlx"
)

type DAO interface {
	Close()
	NamedQuery(query string, filter interface{}) (*sqlx.Rows, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	GetDB() *sqlx.DB
}

type dao struct {
	d *sqlx.DB
}

func NewDao(connectionUri string) DAO {
	db, err := ddsqlx.Connect("postgres", connectionUri)
	if err != nil {
		fmt.Println("Connection error: ", err.Error())
		return nil
	}
	return &dao{
		d: db,
	}
}

func (d *dao) Close() {
	d.d.Close()
}

func (d *dao) NamedQuery(query string, filter interface{}) (*sqlx.Rows, error) {
	return d.d.NamedQuery(query, filter)
}

func (d *dao) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return d.d.NamedExec(query, arg)
}

func (d *dao) GetDB() *sqlx.DB {
	return d.d
}
