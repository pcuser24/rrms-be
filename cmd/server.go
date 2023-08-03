package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	bjwt2 "github.com/user2410/rrms-backend/internal/bjwt"
	"github.com/user2410/rrms-backend/internal/domain/user"
	"github.com/user2410/rrms-backend/internal/infrastructure/http"
	"github.com/user2410/rrms-backend/internal/infrastructure/repositories/db"
	"log"
)

type serverConfig struct {
	DatabaseURL string `mapstructure:"db_url"`
}

type serverCommand struct {
	*cobra.Command
	config *serverConfig
}

func NewServerCommand() *serverCommand {
	c := &serverCommand{}
	c.Command = &cobra.Command{
		Use:   "serve",
		Short: fmt.Sprintf("Http serve for %s", ReadableName),
		Long: fmt.Sprintf(`%s
Manage the APIs for %s from the command line`, Art(), ReadableName),
		Run: c.run,
	}
	c.config = newServerConfig(c.Command)
	return c
}

func (c *serverCommand) run(cmd *cobra.Command, args []string) {
	server := http.NewServer()
	dao := db.NewDao(c.config.DatabaseURL)
	if dao == nil {
		log.Println("Error while initializing database connection")
		return
	}
	defer dao.Close()

	bjwt := bjwt2.NewBjwt("DchinhSecretKey")

	userRepo := user.NewUserRepo(dao)
	userService := user.NewUserService(userRepo)
	userAdapter := user.NewAdapter(userService, bjwt)
	userAdapter.RegisterServer(server)
	server.Start()
}

func newServerConfig(cmd *cobra.Command) *serverConfig {
	//TODO parse env vars or flags
	c := &serverConfig{
		DatabaseURL: "postgres://dchinhpl:dchinhpl@127.0.0.1:32755/rrms?sslmode=disable", //TODO scan from env vars or command flags... *cmd.Flags().String("db_url", "", "URI to connect database"),
	}
	return c
}
