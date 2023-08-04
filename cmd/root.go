/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

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
		Short:         fmt.Sprintf(shortFormat, ReadableName),
		Long:          fmt.Sprintf(longFormat, Art(), ReadableName),
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
		NewVersionCommand().Command,
		NewMigrateCommand().Command,
		NewServerCommand().Command,
	)

	c, err := root.Command.ExecuteC()
	if err != nil {
		c.Println(Art())
		c.Println(c.UsageString())
		c.PrintErrf("ERROR: %v\n", err)
		os.Exit(1)
	}
}
