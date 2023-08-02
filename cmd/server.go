package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

type serverCommand struct {
	*cobra.Command
}

func NewServerCommand() *migrateCommand {
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

func (c *serverCommand) run(cmd *cobra.Command, args []string) {
	c.Help()
}
