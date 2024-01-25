package database

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator interface {
	Upgrade() (bool, error)
	Downgrade() (bool, error)
	Status() (uint, bool, error)
}

type migrator struct {
	m *migrate.Migrate
}

func NewMigrator(migrationDir, dbUrl string) (Migrator, error) {
	migration, err := migrate.New(fmt.Sprintf("file://%s", migrationDir), dbUrl)
	if err != nil {
		return nil, err
	}

	return &migrator{
		m: migration,
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
