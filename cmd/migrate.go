package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/user2410/rrms-backend/internal/infrastructure/repositories/db"
	"log"
	"os"
	"path/filepath"
)

type migrateConfig struct {
	DatabaseURL          string `mapstructure:"db_url"`
	DatabaseMigrationDir string `mapstructure:"db_migration_dir"`
}

type migrateCommand struct {
	*cobra.Command
}

func NewMigrateCommand() *migrateCommand {
	c := &migrateCommand{}
	c.Command = &cobra.Command{
		Use:   "migrate",
		Short: fmt.Sprintf("Database migration manager for %s", ReadableName),
		Long: fmt.Sprintf(`%s
Manage the database migrations for %s from the command line`, Art(), ReadableName),
		Run: c.run,
	}
	c.Command.AddCommand(
		newMigrateUpCommand().Command,
	)
	return c
}

func (c *migrateCommand) run(cmd *cobra.Command, args []string) {
	c.Help()
}

type migrateUpCommand struct {
	*cobra.Command
	config *migrateConfig
}

func newMigrateUpCommand() *migrateUpCommand {
	c := &migrateUpCommand{}
	c.Command = &cobra.Command{
		Use:   "up",
		Short: "Upgrade a migration on the permission API database",
		Run:   c.run,
	}
	c.config = newMigrationConfig(c.Command)
	return c
}

func (c *migrateUpCommand) run(cmd *cobra.Command, args []string) {
	m, err := initMigrator(c.config)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	res, err := m.Upgrade()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	if !res {
		log.Println("The database is already up to date")
	}
	log.Println("The database has been upgraded successfully")
}

func initMigrator(conf *migrateConfig) (db.Migrator, error) {
	dao := db.NewDao(conf.DatabaseURL)
	mDir, err := filepath.Abs(conf.DatabaseMigrationDir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve the migration directory: %w", err)
	}
	m, err := db.NewMigrator(dao, mDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create the migration manager: %w", err)
	}
	return m, nil
}

func newMigrationConfig(cmd *cobra.Command) *migrateConfig {
	c := &migrateConfig{
		DatabaseURL:          "postgres://dchinhpl:dchinhpl@127.0.0.1:32755/rrms?sslmode=disable", //TODO scan from env vars or command flags... *cmd.Flags().String("db_url", "", "URI to connect database"),
		DatabaseMigrationDir: "./migrations",                                                      //TODO scan from env vars or command flags... *cmd.Flags().String("db_migration_dir", "", "Migration Directory"),
	}
	return c
}
