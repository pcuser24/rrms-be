package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type migrateSeedCommand struct {
	*cobra.Command
	config *migrateConfig
}

func NewMigrateSeedCommand() *migrateSeedCommand {
	c := &migrateSeedCommand{}
	c.Command = &cobra.Command{
		Use:   "seed",
		Short: "Seed the database with initial data",
		Run:   c.run,
	}
	c.config = newMigrationConfig(c.Command)
	return c
}

func (c *migrateSeedCommand) run(cmd *cobra.Command, args []string) {
	dao, err := database.NewPostgresDAO(c.config.DatabaseURL)
	if err != nil {
		log.Fatal("failed to create the database connection", err)
	}
	defer dao.Close()
}
