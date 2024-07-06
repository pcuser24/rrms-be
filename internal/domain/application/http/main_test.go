package http

import (
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/require"
	repos "github.com/user2410/rrms-backend/internal/domain/_repos"
	application "github.com/user2410/rrms-backend/internal/domain/application/service"
	listing_service "github.com/user2410/rrms-backend/internal/domain/listing/service"
	misc_service "github.com/user2410/rrms-backend/internal/domain/misc/service"
	"github.com/user2410/rrms-backend/internal/domain/reminder"
	"go.uber.org/mock/gomock"

	"github.com/user2410/rrms-backend/internal/infrastructure/aws/s3"
	"github.com/user2410/rrms-backend/internal/infrastructure/http"
	"github.com/user2410/rrms-backend/internal/utils/random"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

type server struct {
	domainRepo repos.DomainRepo
	tokenMaker token.Maker
	router     http.Server
}

func newTestServer(
	t *testing.T,
	mockCtrl *gomock.Controller,
) *server {

	tokenMaker, err := token.NewJWTMaker(random.RandomAlphanumericStr(32))
	require.NoError(t, err)
	require.NotNil(t, tokenMaker)

	// initialize services
	domainRepo := repos.NewDomainRepoFromMockCtrl(mockCtrl)

	reminderService := reminder.NewService(domainRepo)
	// TODO: mock notification endpoint
	miscService := misc_service.NewService(domainRepo, nil, cron.New())
	// TODO: mock s3 client
	s3Client := s3.NewMockS3Client(mockCtrl)
	applicationService := application.NewService(domainRepo, reminderService, miscService, s3Client, "", nil, "https://rrms.rental.vn/")
	lService := listing_service.NewService(domainRepo, "", nil) // NOTE: leave esClient nil for now

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
	NewAdapter(lService, applicationService).RegisterServer(httpServer.GetApiRoute(), tokenMaker)

	return &server{
		domainRepo: domainRepo,
		tokenMaker: tokenMaker,
		router:     httpServer,
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
