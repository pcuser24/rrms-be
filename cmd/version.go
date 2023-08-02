package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

var version = "development"
var goVersion = fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
var buildStamp = ""

const (
	codeName     = "rrms"
	ReadableName = "RRMS"
	figlet       = `
RRRRRRRRRRRRRRRRRRRR         RRRRRRRRRRRRRRRRRRRR         MMMMMMMMM              MMMMMMMM      		SSSSSSSSSSSSSSS
R::::::R     R:::::R         R::::::R     R:::::R         M:::::::M             M:::::::M     	SS:::::::::::::::S
R::::::R     R:::::R         R::::::R     R:::::R         M::::::::M           M::::::::M       S:::::SSSSSS::::::S
RR:::::R     R:::::R         RR:::::R     R:::::R         M:::::::::M         M:::::::::M      S:::::S     SSSSSSS
R::::R       R:::::R           R::::R     R:::::R         M::::::::::M       M::::::::::M      S:::::S
  R::::R     R:::::R           R::::R     R:::::R         M:::::::::::M     M:::::::::::M      S:::::S
  R::::RRRRRR:::::R            R::::RRRRRR:::::R          M:::::::M::::M   M::::M:::::::M      S::::SSSS
  R:::::::::::::RR             R:::::::::::::RR           M::::::M M::::M M::::M M::::::M       SS::::::SSSSS
  R::::RRRRRR:::::R            R::::RRRRRR:::::R          M::::::M  M::::M::::M  M::::::M         SSS::::::::SS
  R::::R     R:::::R           R::::R     R:::::R         M::::::M   M:::::::M   M::::::M            SSSSSS::::S
  R::::R     R:::::R           R::::R     R:::::R         M::::::M    M:::::M    M::::::M                 S:::::S
  R::::R     R:::::R           R::::R     R:::::R         M::::::M     MMMMM     M::::::M                 S:::::S
RR:::::R     R:::::R         RR:::::R     R:::::R         M::::::M               M::::::M       SSSSSSS     S:::::S
R::::::R     R:::::R         R::::::R     R:::::R         M::::::M               M::::::M         S::::::SSSSSS:::::S
R::::::R     R:::::R         R::::::R     R:::::R         M::::::M               M::::::M           S:::::::::::::::SS
RRRRRRRR     RRRRRRR         RRRRRRRR     RRRRRRR         MMMMMMMM               MMMMMMMM             SSSSSSSSSSSSSSS

                                                        %s %s
`
	versionShortFormat = "Print the Version Information"
	versionLongFormat  = `%s
Print the Version information`
)

type versionCommand struct {
	*cobra.Command
}

func NewVersionCommand() *versionCommand {
	vc := &versionCommand{}
	vc.Command = &cobra.Command{
		Use:   "version",
		Short: versionShortFormat,
		Long:  fmt.Sprintf(versionLongFormat, Art()),
		Run:   vc.run,
	}
	return vc
}

func (vc *versionCommand) run(cmd *cobra.Command, args []string) {
	cmd.Println(Art())
	cmd.Println(codeName)
	cmd.Println(fmt.Sprintf("  version: %s", version))
	cmd.Println(fmt.Sprintf("  go: %s", goVersion))
	cmd.Println(fmt.Sprintf("  built at: %s", buildStamp))
}

func Art() string {
	return fmt.Sprintf(figlet, ReadableName, version)
}
