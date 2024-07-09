package http

import (
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/stretchr/testify/require"
	repos "github.com/user2410/rrms-backend/internal/domain/_repos"
	auth_service "github.com/user2410/rrms-backend/internal/domain/auth/service"
	property_service "github.com/user2410/rrms-backend/internal/domain/property/service"
	"github.com/user2410/rrms-backend/internal/infrastructure/aws/s3"
	"github.com/user2410/rrms-backend/internal/infrastructure/http"
	"github.com/user2410/rrms-backend/internal/utils/random"
	"github.com/user2410/rrms-backend/internal/utils/token"
	"go.uber.org/mock/gomock"
)

type server struct {
	domainRepo repos.DomainRepo
	tokenMaker token.Maker
	router     http.Server
}

func newTestServer(t *testing.T, ctrl *gomock.Controller) *server {

	tokenMaker, err := token.NewJWTMaker(random.RandomAlphanumericStr(32))
	require.NoError(t, err)
	require.NotNil(t, tokenMaker)

	domainRepo := repos.NewDomainRepoFromMockCtrl(ctrl)

	s3Client := s3.NewMockS3Client(ctrl)

	// initialize service
	authService := auth_service.NewService(domainRepo, tokenMaker, time.Hour, time.Hour)
	service := property_service.NewService(domainRepo, s3Client, "", nil, nil, nil)

	// initialize http router
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
	NewAdapter(service).RegisterServer(httpServer.GetApiRoute(), tokenMaker, authService)

	return &server{
		domainRepo: domainRepo,
		tokenMaker: tokenMaker,
		router:     httpServer,
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
