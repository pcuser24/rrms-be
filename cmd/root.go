/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/user2410/rrms-backend/cmd/migrate"
	"github.com/user2410/rrms-backend/cmd/seed"
	"github.com/user2410/rrms-backend/cmd/server"
	"github.com/user2410/rrms-backend/cmd/version"

	"github.com/spf13/cobra"
)

type rootCommand struct {
	*cobra.Command
}

func (rc *rootCommand) run(c *cobra.Command, args []string) {
	fmt.Println("Root command 1234")
	c.Help()
}

const shortFormat string = "Command Line Interface Manager for %s"
const longFormat string = `%s
Manage %s from the command line`

func newRootCommand() *rootCommand {
	rc := &rootCommand{}
	rc.Command = &cobra.Command{
		Use:           "rrmsd",
		Short:         fmt.Sprintf(shortFormat, version.ReadableName),
		Long:          fmt.Sprintf(longFormat, version.Art(), version.ReadableName),
		Run:           rc.run,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	return rc
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	root := newRootCommand()
	root.Command.AddCommand(
		migrate.NewMigrateCommand().Command,
		seed.NewSeedCommand().Command,
		server.NewServerCommand().Command,
		version.NewVersionCommand().Command,
	)

	c, err := root.Command.ExecuteC()
	if err != nil {
		c.Println(version.Art())
		c.Println(c.UsageString())
		c.PrintErrf("ERROR: %v\n", err)
		os.Exit(1)
	}
}
