package db

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

type Migrator interface {
	Upgrade() (bool, error)
	Downgrade() (bool, error)
	Status() (uint, bool, error)
}

type migrator struct {
	m *migrate.Migrate
}

func NewMigrator(conn DAO, migrationDir string) (Migrator, error) {
	if conn == nil {
		return nil, fmt.Errorf("the database connection passed in is invalid")
	}
	dbConf := &postgres.Config{}
	dbIns, err := postgres.WithInstance(conn.GetDB().DB, dbConf)
	if err != nil {
		return nil, fmt.Errorf("failed to create a database instance from the connection: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationDir),
		"postgres",
		dbIns,
	)
	if err != nil {
		return nil, err
	}
	return &migrator{
		m: m,
	}, err
}

//func NewMigratorFromDB(db *sql.DB, migrationDir string) (Migrator, error) {
//	conn := NewConnectionFromDB(db)
//	return NewMigrator(conn, migrationDir)
//}

func (m *migrator) Upgrade() (bool, error) {
	var updateOccurred bool
	err := m.m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		err = nil
	} else if err == nil {
		updateOccurred = true
	}
	return updateOccurred, err
}
func (m *migrator) Downgrade() (bool, error) {
	var downgradeOccurred bool
	_, _, err := m.m.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		return downgradeOccurred, nil
	}
	err = m.m.Steps(-1)
	if err == nil {
		downgradeOccurred = true
	}
	return downgradeOccurred, err
}

func (m *migrator) Status() (uint, bool, error) {
	version, dirty, err := m.m.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		err = nil
	}
	return version, dirty, err
}
