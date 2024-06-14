package http

import (
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/stretchr/testify/require"
	application_repo "github.com/user2410/rrms-backend/internal/domain/application/repo"
	application "github.com/user2410/rrms-backend/internal/domain/application/service"
	auth_repo "github.com/user2410/rrms-backend/internal/domain/auth/repo"
	chat_repo "github.com/user2410/rrms-backend/internal/domain/chat/repo"
	listing_repo "github.com/user2410/rrms-backend/internal/domain/listing/repo"
	listing_service "github.com/user2410/rrms-backend/internal/domain/listing/service"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	"github.com/user2410/rrms-backend/internal/domain/reminder"
	reminder_repo "github.com/user2410/rrms-backend/internal/domain/reminder/repo"
	unit_repo "github.com/user2410/rrms-backend/internal/domain/unit/repo"
	"go.uber.org/mock/gomock"

	"github.com/user2410/rrms-backend/internal/infrastructure/aws/s3"
	"github.com/user2410/rrms-backend/internal/infrastructure/http"
	"github.com/user2410/rrms-backend/internal/utils/random"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

type server struct {
	ur         unit_repo.Repo
	pr         property_repo.Repo
	lr         listing_repo.Repo
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
	propertyRepo := property_repo.NewMockRepo(mockCtrl)
	unitRepo := unit_repo.NewMockRepo(mockCtrl)
	listingRepo := listing_repo.NewMockRepo(mockCtrl)
	applicationRepo := application_repo.NewMockRepo(mockCtrl)
	reminderRepo := reminder_repo.NewMockRepo(mockCtrl)
	authRepo := auth_repo.NewMockRepo(mockCtrl)
	chatRepo := chat_repo.NewMockRepo(mockCtrl)

	s3Client := s3.NewMockS3Client(mockCtrl)

	rService := reminder.NewService(reminderRepo)
	// miscService := misc.NewService()
	aService := application.NewService(applicationRepo, authRepo, chatRepo, listingRepo, propertyRepo, unitRepo, rService, nil, s3Client, "", "https://rrms.rental.vn/")
	lService := listing_service.NewService(listingRepo, propertyRepo, unitRepo, nil, "") // NOTE: leave paymentRepo nil for now

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
	NewAdapter(lService, aService).RegisterServer(httpServer.GetApiRoute(), tokenMaker)

	return &server{
		pr:         propertyRepo,
		ur:         unitRepo,
		lr:         listingRepo,
		tokenMaker: tokenMaker,
		router:     httpServer,
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
