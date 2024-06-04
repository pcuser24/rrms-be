package http

import (
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/stretchr/testify/require"
	"github.com/user2410/rrms-backend/internal/domain/application/repo"
	application "github.com/user2410/rrms-backend/internal/domain/application/service"
	chat_repo "github.com/user2410/rrms-backend/internal/domain/chat/repo"
	listing_repo "github.com/user2410/rrms-backend/internal/domain/listing/repo"
	listing_service "github.com/user2410/rrms-backend/internal/domain/listing/service"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	"github.com/user2410/rrms-backend/internal/domain/reminder"
	reminder_repo "github.com/user2410/rrms-backend/internal/domain/reminder/repo"
	unit_repo "github.com/user2410/rrms-backend/internal/domain/unit/repo"

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
	applicationRepo repo.Repo, propertyRepo property_repo.Repo, unitRepo unit_repo.Repo, listingRepo listing_repo.Repo, chatRepo chat_repo.Repo, reminderRepo reminder_repo.Repo,
) *server {

	tokenMaker, err := token.NewJWTMaker(random.RandomAlphanumericStr(32))
	require.NoError(t, err)
	require.NotNil(t, tokenMaker)

	// initialize services
	rService := reminder.NewService(reminderRepo)
	aService := application.NewService(applicationRepo, chatRepo, listingRepo, propertyRepo, rService, nil, "")
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
