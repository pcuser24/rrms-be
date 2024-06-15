package http

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"
	repos "github.com/user2410/rrms-backend/internal/domain/_repos"
	listing_service "github.com/user2410/rrms-backend/internal/domain/listing/service"
	"github.com/user2410/rrms-backend/internal/domain/payment/service/vnpay"
	"github.com/user2410/rrms-backend/internal/infrastructure/http"
	"github.com/user2410/rrms-backend/internal/utils"
	"go.uber.org/mock/gomock"
)

type ServerConfig struct {
	// DatabaseURL string `mapstructure:"DB_URL" validate:"required,uri"`
	VnpTmnCode    string `mapstructure:"VNP_TMNCODE" validate:"required"`
	VnpHashSecret string `mapstructure:"VNP_HASHSECRET" validate:"required"`
	VnpUrl        string `mapstructure:"VNP_URL" validate:"required"`
	VnpApi        string `mapstructure:"VNP_API" validate:"required"`
	VnpReturnUrl  string `mapstructure:"VNP_RETURNURL" validate:"required"`
}

var (
	basePath = utils.GetBasePath()
	conf     ServerConfig
)

type server struct {
	router http.Server
}

func newTestServer(t *testing.T, ctrl *gomock.Controller) *server {

	domainRepo := repos.NewDomainRepoFromMockCtrl(ctrl)
	listingService := listing_service.NewService(domainRepo, "")
	vnpService := vnpay.NewVnpayService(domainRepo, listingService, conf.VnpTmnCode, conf.VnpHashSecret, conf.VnpUrl, conf.VnpApi)

	httpServer := http.NewServer(
		fiber.Config{
			ReadTimeout:  1 * time.Second,
			WriteTimeout: 1 * time.Second,
		},
		cors.Config{
			AllowOrigins: "*",
			AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		},
	)

	NewAdapter(vnpService).RegisterServer(httpServer.GetApiRoute(), nil)

	return &server{
		router: httpServer,
	}
}

func TestMain(m *testing.M) {
	viper.AddConfigPath(basePath)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("failed to read config file: %v", err)
	}
	err = viper.Unmarshal(&conf)
	if err != nil {
		log.Printf("failed to unmarshal config file: %v", err)
	}
	v := validator.New()
	err = v.Struct(&conf)
	if err != nil {
		log.Printf("invalid or missing fields in config file: %v", err)
	}

	os.Exit(m.Run())
}
