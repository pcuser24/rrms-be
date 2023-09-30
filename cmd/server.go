package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/user2410/rrms-backend/internal/domain/auth"
	"github.com/user2410/rrms-backend/internal/domain/listing"
	"github.com/user2410/rrms-backend/internal/domain/property"
	"github.com/user2410/rrms-backend/internal/domain/rental"
	"github.com/user2410/rrms-backend/internal/domain/storage"
	"github.com/user2410/rrms-backend/internal/domain/unit"
	"github.com/user2410/rrms-backend/internal/infrastructure/aws/s3"
	db "github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/infrastructure/http"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

type ServerConfig struct {
	DatabaseURL string `mapstructure:"DB_URL" validate:"required,uri"`

	AllowOrigins string `mapstructure:"ALLOW_ORIGINS" validate:"required"`

	TokenMaker      string        `mapstructure:"TOKEN_MAKER" validate:"required"`
	TokenSecreteKey string        `mapstructure:"TOKEN_SECRET_KEY" validate:"required"`
	AccessTokenTTL  time.Duration `mapstructure:"ACCESS_TOKEN_TTL" validate:"required"`

	AWSRegion          string `mapstructure:"AWS_REGION" validate:"required"`
	AWSAccessKeyID     string `mapstructure:"AWS_ACCESS_KEY_ID" validate:"required"`
	AWSSecretAccessKey string `mapstructure:"AWS_SECRET_ACCESS_KEY" validate:"required"`
	AWSS3BucketName    string `mapstructure:"AWS_S3_BUCKET_NAME" validate:"required"`
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
	// setup database
	dao, err := db.NewDAO(c.config.DatabaseURL)
	if err != nil {
		log.Fatal("Error while initializing database connection: ", err)
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

	// setup S3 client
	s3Storage, err := s3.NewAWSS3StorageService(
		c.config.AWSRegion,
		c.config.AWSAccessKeyID,
		c.config.AWSSecretAccessKey,
		c.config.AWSS3BucketName,
	)
	if err != nil {
		log.Fatal("Error while initializing AWS S3 client", err)
	}

	// setup http server
	c.httpServer = http.NewServer(
		fiber.Config{
			ReadTimeout:  1 * time.Second,
			WriteTimeout: 1 * time.Second,
		},
		cors.Config{
			AllowOrigins: c.config.AllowOrigins,
			AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		},
	)
	apiRoute := c.httpServer.GetApiRoute()

	authRepo := auth.NewUserRepo(dao)
	authService := auth.NewUserService(authRepo, tokenMaker, c.config.AccessTokenTTL)
	auth.NewAdapter(authService).RegisterServer(apiRoute)
	propertyRepo := property.NewRepo(dao)
	propertyService := property.NewService(propertyRepo)
	property.NewAdapter(propertyService).RegisterServer(apiRoute, tokenMaker)
	unitRepo := unit.NewRepo(dao)
	unitService := unit.NewService(unitRepo)
	unit.NewAdapter(unitService, propertyService).RegisterServer(apiRoute, tokenMaker)
	listingRepo := listing.NewRepo(dao)
	listingService := listing.NewService(listingRepo)
	listing.NewAdapter(listingService, propertyService, unitService).RegisterServer(apiRoute, tokenMaker)
	rentalRepo := rental.NewRepo(dao)
	rentalService := rental.NewService(rentalRepo)
	rental.NewAdapter(rentalService).RegisterServer(apiRoute)

	storageService := storage.NewService(s3Storage)
	storage.NewAdapter(storageService).RegisterServer(apiRoute, tokenMaker)
}

func (c *serverCommand) run(cmd *cobra.Command, args []string) {
	c.setup(cmd, args)

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
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

	os.Exit(0)
}
