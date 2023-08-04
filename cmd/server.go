package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/user2410/rrms-backend/internal/domain/user"
	db "github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/infrastructure/http"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

type ServerConfig struct {
	DatabaseURL string `mapstructure:"DB_URL" validate:"required,uri"`

	TokenMaker      string        `mapstructure:"TOKEN_MAKER" validate:"required"`
	TokenSecreteKey string        `mapstructure:"TOKEN_SECRET_KEY" validate:"required"`
	AccessTokenTTL  time.Duration `mapstructure:"ACCESS_TOKEN_TTL" validate:"required"`
}

type serverCommand struct {
	*cobra.Command
	config     *ServerConfig
	dao        db.DAO
	httpServer http.Server
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

func newServerConfig(cmd *cobra.Command) *ServerConfig {
	//TODO parse env vars or flags
	var conf ServerConfig
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Failed to read config file:", err)
	}
	err = viper.Unmarshal(&conf)
	if err != nil {
		log.Fatal("Failed to unmarshal config file:", err)
	}

	v := validator.New()
	err = v.Struct(&conf)
	if err != nil {
		log.Fatal("Invalid or missing fields in config file: ", err)
	}

	return &conf
}

func (c *serverCommand) setup(cmd *cobra.Command, args []string) {
	var err error

	// setup database
	dao := db.NewDao(c.config.DatabaseURL)
	if dao == nil {
		log.Println("Error while initializing database connection")
		return
	}
	c.dao = dao

	// setup token maker
	var tokenMaker token.Maker
	switch strings.ToUpper(c.config.DatabaseURL) {
	case "PASETO":
		tokenMaker, err = token.NewPasetoMaker(c.config.TokenSecreteKey)
		if err != nil {
			log.Fatal(err)
		}
	default:
		tokenMaker, err = token.NewJWTMaker(c.config.TokenSecreteKey)
		if err != nil {
			log.Fatal(err)
		}
	}

	// setup http server
	c.httpServer = http.NewServer(fiber.Config{
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	})
	userRepo := user.NewUserRepo(dao)
	userService := user.NewUserService(userRepo, tokenMaker, c.config.AccessTokenTTL)
	userAdapter := user.NewAdapter(userService)
	userAdapter.RegisterServer(c.httpServer.GetFibApp())
}

func (c *serverCommand) run(cmd *cobra.Command, args []string) {
	c.setup(cmd, args)

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt)
	go func() {
		<-exitCh
		fmt.Println("Gracefully shutting down...")
		err := c.httpServer.Shutdown()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if err := c.httpServer.Start(8000); err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Running cleanup tasks...")
	c.shutdown()
}

func (c *serverCommand) shutdown() {
	c.dao.Close()

	if err := c.httpServer.Shutdown(); err != nil {
		log.Fatal(err)
	}
}
