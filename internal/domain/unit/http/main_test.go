package http

import (
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/stretchr/testify/require"
	application_repo "github.com/user2410/rrms-backend/internal/domain/application/repo"
	auth_repo "github.com/user2410/rrms-backend/internal/domain/auth/repo"
	listing_repo "github.com/user2410/rrms-backend/internal/domain/listing/repo"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	property_service "github.com/user2410/rrms-backend/internal/domain/property/service"
	rental_repo "github.com/user2410/rrms-backend/internal/domain/rental/repo"
	"github.com/user2410/rrms-backend/internal/domain/unit"
	"github.com/user2410/rrms-backend/internal/domain/unit/repo"
	"github.com/user2410/rrms-backend/internal/infrastructure/http"
	"github.com/user2410/rrms-backend/internal/utils/random"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

type server struct {
	ur         repo.Repo
	pr         property_repo.Repo
	tokenMaker token.Maker
	router     http.Server
}

func newTestServer(t *testing.T, propertyRepo property_repo.Repo, unitRepo repo.Repo, listingRepo listing_repo.Repo, applicationRepo application_repo.Repo, rentalRepo rental_repo.Repo, authRepo auth_repo.Repo) *server {

	tokenMaker, err := token.NewJWTMaker(random.RandomAlphanumericStr(32))
	require.NoError(t, err)
	require.NotNil(t, tokenMaker)

	// initialize uService
	uService := unit.NewService(unitRepo, nil, "")
	pService := property_service.NewService(propertyRepo, unitRepo, listingRepo, applicationRepo, rentalRepo, authRepo, nil, "")

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
	NewAdapter(uService, pService).RegisterServer(httpServer.GetApiRoute(), tokenMaker)

	return &server{
		ur:         unitRepo,
		pr:         propertyRepo,
		tokenMaker: tokenMaker,
		router:     httpServer,
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
