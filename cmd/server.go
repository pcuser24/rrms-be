package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	application_asynctask "github.com/user2410/rrms-backend/internal/domain/application/asynctask"
	auth_asynctask "github.com/user2410/rrms-backend/internal/domain/auth/asynctask"

	application_repo "github.com/user2410/rrms-backend/internal/domain/application/repo"
	auth_repo "github.com/user2410/rrms-backend/internal/domain/auth/repo"
	listing_repo "github.com/user2410/rrms-backend/internal/domain/listing/repo"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	unit_repo "github.com/user2410/rrms-backend/internal/domain/unit/repo"

	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	property_http "github.com/user2410/rrms-backend/internal/domain/property/http"
	unit_http "github.com/user2410/rrms-backend/internal/domain/unit/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/hibiken/asynq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/user2410/rrms-backend/internal/domain/application"
	"github.com/user2410/rrms-backend/internal/domain/auth"
	"github.com/user2410/rrms-backend/internal/domain/listing"
	"github.com/user2410/rrms-backend/internal/domain/property"
	"github.com/user2410/rrms-backend/internal/domain/rental"
	"github.com/user2410/rrms-backend/internal/domain/storage"
	"github.com/user2410/rrms-backend/internal/domain/unit"
	"github.com/user2410/rrms-backend/internal/infrastructure/asynctask"
	"github.com/user2410/rrms-backend/internal/infrastructure/aws/s3"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/infrastructure/email"
	"github.com/user2410/rrms-backend/internal/infrastructure/http"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

type ServerConfig struct {
	DatabaseURL string `mapstructure:"DB_URL" validate:"required,uri"`

	AllowOrigins string `mapstructure:"ALLOW_ORIGINS" validate:"required"`

	TokenMaker      string        `mapstructure:"TOKEN_MAKER" validate:"required"`
	TokenSecreteKey string        `mapstructure:"TOKEN_SECRET_KEY" validate:"required"`
	AccessTokenTTL  time.Duration `mapstructure:"ACCESS_TOKEN_TTL" validate:"required"`
	RefreshTokenTTL time.Duration `mapstructure:"REFRESH_TOKEN_TTL" validate:"required"`

	AWSRegion          string  `mapstructure:"AWS_REGION" validate:"required"`
	AWSAccessKeyID     string  `mapstructure:"AWS_ACCESS_KEY_ID" validate:"required"`
	AWSSecretAccessKey string  `mapstructure:"AWS_SECRET_ACCESS_KEY" validate:"required"`
	AWSS3Endpoint      *string `mapstructure:"AWS_S3_ENDPOINT" validate:"omitempty"`
	AWSS3ImageBucket   string  `mapstructure:"AWS_S3_IMAGE_BUCKET" validate:"required"`

	EmailSenderName     string `mapstructure:"EMAIL_SENDER_NAME" validate:"required"`
	EmailSenderAddress  string `mapstructure:"EMAIL_SENDER_ADDRESS" validate:"required"`
	EmailSenderPassword string `mapstructure:"EMAIL_SENDER_PASSWORD" validate:"required"`

	AsynqRedisAddress string `mapstructure:"ASYNQ_REDIS_ADDRESS" validate:"required"`
}

type internalServices struct {
	AuthService        auth.Service
	PropertyService    property.Service
	UnitService        unit.Service
	ListingService     listing.Service
	RentalService      rental.Service
	ApplicationService application.Service
	StorageService     storage.Service
}

type serverCommand struct {
	*cobra.Command
	tokenMaker           token.Maker
	emailSender          email.EmailSender
	config               *ServerConfig
	dao                  database.DAO
	internalServices     internalServices
	httpServer           http.Server
	asyncTaskDistributor asynctask.Distributor
	asyncTaskProcessor   asynctask.Processor
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

func (c *serverCommand) run(cmd *cobra.Command, args []string) {
	c.setup(cmd, args)

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	go func() {
		<-exitCh
		fmt.Println("Gracefully shutting down...")
		c.shutdown()
	}()

	go c.runAsyncTaskProcessor()
	c.runHttpServer()

}

func (c *serverCommand) shutdown() {
	log.Println("Running cleanup tasks...")

	c.dao.Close()

	if err := c.httpServer.Shutdown(); err != nil {
		log.Fatal(err)
	}

	c.asyncTaskProcessor.Shutdown()

	os.Exit(0)
}

/* -------------------------------------------------------------------------- */
/*                       setups components of the server                      */
/* -------------------------------------------------------------------------- */

func (c *serverCommand) setup(cmd *cobra.Command, args []string) {
	// setup database
	dao, err := database.NewPostgresDAO(c.config.DatabaseURL)
	if err != nil {
		log.Fatal("Error while initializing database connection: ", err)
	}
	c.dao = dao

	// setup token maker
	if strings.ToUpper(c.config.TokenMaker) == "PASETO" {
		c.tokenMaker, err = token.NewPasetoMaker(c.config.TokenSecreteKey)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		c.tokenMaker, err = token.NewJWTMaker(c.config.TokenSecreteKey)
		if err != nil {
			log.Fatal(err)
		}
	}

	// setup mailer
	c.emailSender = email.NewGmailSender(
		c.config.EmailSenderName,
		c.config.EmailSenderAddress,
		c.config.EmailSenderPassword,
	)

	// setup S3 client
	s3Client, err := s3.NewS3Client(
		c.config.AWSRegion,
		c.config.AWSAccessKeyID,
		c.config.AWSSecretAccessKey,
		c.config.AWSS3Endpoint,
	)
	if err != nil {
		log.Fatal("Error while initializing AWS S3 client", err)
	}

	// setup asynq task distributor and processor
	c.setupAsyncTaskProcessor(c.emailSender)

	// setup internal services
	c.setupInternalServices(
		dao,
		s3Client,
	)

	// setup http server
	c.setupHttpServer()
}

func (c *serverCommand) setupInternalServices(
	dao database.DAO,
	s3Client *s3.S3Client,
) {
	c.asyncTaskDistributor = asynctask.NewRedisTaskDistributor(asynq.RedisClientOpt{
		Addr: c.config.AsynqRedisAddress,
	})

	authRepo := auth_repo.NewRepo(dao)
	authTaskDistributor := auth_asynctask.NewTaskDistributor(c.asyncTaskDistributor)
	c.internalServices.AuthService = auth.NewService(
		authRepo,
		c.tokenMaker, c.config.AccessTokenTTL, c.config.RefreshTokenTTL,
		authTaskDistributor,
	)
	propertyRepo := property_repo.NewRepo(dao)
	unitRepo := unit_repo.NewRepo(dao)
	listingRepo := listing_repo.NewRepo(dao)
	rentalRepo := rental.NewRepo(dao)
	applicationRepo := application_repo.NewRepo(dao)

	s := storage.NewStorage(s3Client, c.config.AWSS3ImageBucket)

	c.internalServices.PropertyService = property.NewService(propertyRepo, unitRepo)
	c.internalServices.UnitService = unit.NewService(unitRepo)
	c.internalServices.ListingService = listing.NewService(listingRepo)
	c.internalServices.RentalService = rental.NewService(rentalRepo)
	applicationTaskDistributor := application_asynctask.NewTaskDistributor(c.asyncTaskDistributor)
	c.internalServices.ApplicationService = application.NewService(
		applicationRepo,
		applicationTaskDistributor,
	)
	c.internalServices.StorageService = storage.NewService(s)
}

func (c *serverCommand) setupAsyncTaskProcessor(
	mailer email.EmailSender,
) {
	c.asyncTaskProcessor = asynctask.NewRedisTaskProcessor(asynq.RedisClientOpt{
		Addr: c.config.AsynqRedisAddress,
	})

	auth_asynctask.NewTaskProcessor(c.asyncTaskProcessor, mailer).RegisterProcessor()
	application_asynctask.NewTaskProcessor(c.asyncTaskProcessor, mailer).RegisterProcessor()
}

func (c *serverCommand) setupHttpServer() {
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

	auth_http.
		NewAdapter(c.internalServices.AuthService).
		RegisterServer(apiRoute, c.tokenMaker)
	property_http.
		NewAdapter(c.internalServices.PropertyService).
		RegisterServer(apiRoute, c.tokenMaker)
	unit_http.NewAdapter(c.internalServices.UnitService, c.internalServices.PropertyService).
		RegisterServer(apiRoute, c.tokenMaker)
	listing.
		NewAdapter(c.internalServices.ListingService, c.internalServices.PropertyService, c.internalServices.UnitService).
		RegisterServer(apiRoute, c.tokenMaker)
	rental.
		NewAdapter(c.internalServices.RentalService).
		RegisterServer(apiRoute)
	application.
		NewAdapter(c.internalServices.ApplicationService).
		RegisterServer(apiRoute, c.tokenMaker)
	storage.
		NewAdapter(c.internalServices.StorageService).
		RegisterServer(apiRoute, c.tokenMaker)
}

/* -------------------------------------------------------------------------- */
/*                        Run components of the server                        */
/* -------------------------------------------------------------------------- */

func (c *serverCommand) runAsyncTaskProcessor() {
	log.Println("Starting async task processor...")
	if err := c.asyncTaskProcessor.Start(); err != nil {
		log.Fatal("Failed to start task processor:", err)
	}
}

func (c *serverCommand) runHttpServer() {
	log.Println("Starting HTTP server...")
	if err := c.httpServer.Start(8000); err != nil {
		log.Fatal("Failed to start HTTP server:", err)
	}
}
