package server

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/user2410/rrms-backend/cmd/version"
	application_service "github.com/user2410/rrms-backend/internal/domain/application/service"
	"github.com/user2410/rrms-backend/internal/domain/auth"
	"github.com/user2410/rrms-backend/internal/domain/chat"
	listing_service "github.com/user2410/rrms-backend/internal/domain/listing/service"
	payment_service "github.com/user2410/rrms-backend/internal/domain/payment/service"
	property_service "github.com/user2410/rrms-backend/internal/domain/property/service"
	"github.com/user2410/rrms-backend/internal/domain/reminder"
	rental_service "github.com/user2410/rrms-backend/internal/domain/rental/service"
	statistic_service "github.com/user2410/rrms-backend/internal/domain/statistic/service"
	"github.com/user2410/rrms-backend/internal/domain/storage"
	"github.com/user2410/rrms-backend/internal/domain/unit"
	"github.com/user2410/rrms-backend/internal/infrastructure/aws/s3"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/infrastructure/es"
	"github.com/user2410/rrms-backend/internal/infrastructure/http"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

type serverConfig struct {
	Port *uint16 `mapstructure:"PORT" validate:"omitempty"`

	DatabaseURL string `mapstructure:"DB_URL" validate:"required,uri"`

	AllowOrigins string `mapstructure:"ALLOW_ORIGINS" validate:"required"`

	TokenMaker      string        `mapstructure:"TOKEN_MAKER" validate:"required"`
	TokenSecreteKey string        `mapstructure:"TOKEN_SECRET_KEY" validate:"required"`
	AccessTokenTTL  time.Duration `mapstructure:"ACCESS_TOKEN_TTL" validate:"required"`
	RefreshTokenTTL time.Duration `mapstructure:"REFRESH_TOKEN_TTL" validate:"required"`

	AWSRegion string `mapstructure:"AWS_REGION" validate:"required"`
	// AWSAccessKeyID     string  `mapstructure:"AWS_ACCESS_KEY_ID" validate:"required"`
	// AWSSecretAccessKey string  `mapstructure:"AWS_SECRET_ACCESS_KEY" validate:"required"`
	AWSS3Endpoint    *string `mapstructure:"AWS_S3_ENDPOINT" validate:"omitempty"`
	AWSS3ImageBucket string  `mapstructure:"AWS_S3_IMAGE_BUCKET" validate:"required"`

	EmailSenderName     string `mapstructure:"EMAIL_SENDER_NAME" validate:"omitempty"`
	EmailSenderAddress  string `mapstructure:"EMAIL_SENDER_ADDRESS" validate:"omitempty"`
	EmailSenderPassword string `mapstructure:"EMAIL_SENDER_PASSWORD" validate:"omitempty"`

	ResendAPIKey string `mapstructure:"RESEND_API_KEY" validate:"omitempty"`

	VnpTmnCode    string `mapstructure:"VNP_TMNCODE" validate:"required"`
	VnpHashSecret string `mapstructure:"VNP_HASHSECRET" validate:"required"`
	VnpUrl        string `mapstructure:"VNP_URL" validate:"required"`
	VnpApi        string `mapstructure:"VNP_API" validate:"required"`

	// Elasticsearch
	ElasticsearchAddresses  *string `mapstructure:"ELASTICSEARCH_ADDRESSES" validate:"omitempty"`
	ElasticsearchUsername   *string `mapstructure:"ELASTICSEARCH_USERNAME" validate:"omitempty"`
	ElasticsearchPassword   *string `mapstructure:"ELASTICSEARCH_PASSWORD" validate:"omitempty"`
	ElasticsearchCACertPath *string `mapstructure:"ELASTICSEARCH_CACERT_PATH" validate:"omitempty"`
	ElasticsearchURL        *string `mapstructure:"ELASTICSEARCH_URL" validate:"omitempty"`
	ElasticsearchCloudID    *string `mapstructure:"ELASTICSEARCH_CLOUD_ID" validate:"omitempty"`
	ElasticsearchAPIKey     *string `mapstructure:"ELASTICSEARCH_API_KEY" validate:"omitempty"`
}

type internalServices struct {
	AuthService        auth.Service
	PropertyService    property_service.Service
	UnitService        unit.Service
	ListingService     listing_service.Service
	RentalService      rental_service.Service
	ApplicationService application_service.Service
	StorageService     storage.Service
	ReminderService    reminder.Service
	PaymentService     payment_service.Service
	ChatService        chat.Service
	StatisticService   statistic_service.Service
}

type serverCommand struct {
	*cobra.Command
	config        *serverConfig
	cronScheduler *cron.Cron
	tokenMaker    token.Maker
	// emailSender      email.EmailSender
	s3Client         *s3.S3Client
	dao              database.DAO
	internalServices internalServices
	httpServer       http.Server
	elasticsearch    *es.ElasticSearchClient
}

func NewServerCommand() *serverCommand {
	c := &serverCommand{}
	c.Command = &cobra.Command{
		Use:   "serve",
		Short: fmt.Sprintf("Http serve for %s", version.ReadableName),
		Long: fmt.Sprintf(`%s
Manage the APIs for %s from the command line`, version.Art(), version.ReadableName),
		Run: c.run,
	}
	c.config = newServerConfig()
	return c
}

func newServerConfig() *serverConfig {
	var conf serverConfig
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
	c.setup()
	defer c.shutdown()

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)

	errChan := make(chan error, 1)
	c.cronScheduler.Start()
	go c.runHttpServer(errChan)

	select {
	case err := <-errChan:
		log.Println("Error while running server: ", err)
	case <-exitCh:
		log.Println("Gracefully shutting down...")
	}

}

func (c *serverCommand) shutdown() {
	log.Println("Running cleanup tasks...")

	c.dao.Close()

	if err := c.httpServer.Shutdown(); err != nil {
		log.Fatal(err)
	}

	c.cronScheduler.Stop()

	os.Exit(0)
}

/* -------------------------------------------------------------------------- */
/*                       setups components of the server                      */
/* -------------------------------------------------------------------------- */

func (c *serverCommand) setup() {
	// setup database
	dao, err := database.NewPostgresDAO(c.config.DatabaseURL)
	if err != nil {
		log.Fatal("Error while initializing database connection: ", err)
	}
	c.dao = dao

	// setup cron scheduler
	c.cronScheduler = cron.New()

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
	// c.emailSender = email.NewResendSender(c.config.ResendAPIKey)

	// setup S3 client
	c.s3Client, err = s3.NewS3Client(c.config.AWSRegion, c.config.AWSS3Endpoint)
	if err != nil {
		log.Fatal("Error while initializing AWS S3 client", err)
	}

	// setup elasticsearch client
	c.elasticsearch, err = es.NewElasticSearchClient(es.ElasticSearchClientParams{
		Addresses:  c.config.ElasticsearchAddresses,
		Username:   c.config.ElasticsearchUsername,
		Password:   c.config.ElasticsearchPassword,
		CACertPath: c.config.ElasticsearchCACertPath,
		Url:        c.config.ElasticsearchURL,
		CloudID:    c.config.ElasticsearchCloudID,
		APIKey:     c.config.ElasticsearchAPIKey,
	})
	if err != nil {
		log.Fatal("Error while initializing Elasticsearch client", err)
	}

	// new http server
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

	// setup internal services
	c.setupInternalServices()

	c.setupHttpServer()
}
