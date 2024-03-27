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
	"github.com/user2410/rrms-backend/internal/domain/payment/repo"
	"github.com/user2410/rrms-backend/internal/domain/payment/service"
	"github.com/user2410/rrms-backend/internal/domain/payment/service/vnpay"
	"github.com/user2410/rrms-backend/internal/infrastructure/http"
	"github.com/user2410/rrms-backend/internal/utils"
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

func newTestServer(t *testing.T, repo repo.Repo) *server {
	vnpService := vnpay.NewVnpayService(repo, nil, conf.VnpTmnCode, conf.VnpHashSecret, conf.VnpUrl, conf.VnpApi)
	paymentService := service.NewService(repo)

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

	NewAdapter(paymentService, vnpService).RegisterServer(httpServer.GetApiRoute(), nil)

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
