package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/user2410/rrms-backend/cmd"
	"os"
)

type rootCommand struct {
	*cobra.Command
}

func (rc *rootCommand) run(c *cobra.Command, args []string) {
	c.Help()
}

const shortFormat string = "Command Line Interface Manager for %s"
const longFormat string = `%s
Manage %s from the command line`

func newRootCommand() *rootCommand {
	rc := &rootCommand{}
	rc.Command = &cobra.Command{
		Use:           "rrmsd",
		Short:         fmt.Sprintf(shortFormat, cmd.ReadableName),
		Long:          fmt.Sprintf(longFormat, cmd.Art(), cmd.ReadableName),
		Run:           rc.run,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	return rc
}

func main() {
	root := newRootCommand()
	root.Command.AddCommand(
		cmd.NewVersionCommand().Command,
		cmd.NewMigrateCommand().Command,
		cmd.NewServerCommand().Command,
	)

	c, err := root.Command.ExecuteC()
	if err != nil {
		c.Println(cmd.Art())
		c.Println(c.UsageString())
		c.PrintErrf("ERROR: %v\n", err)
		os.Exit(1)
	}
}
