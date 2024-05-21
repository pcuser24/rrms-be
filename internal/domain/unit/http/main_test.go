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

func newTestServer(t *testing.T, pr property_repo.Repo, ur repo.Repo, lr listing_repo.Repo, ar application_repo.Repo, rr rental_repo.Repo, authRepo auth_repo.Repo) *server {

	tokenMaker, err := token.NewJWTMaker(random.RandomAlphanumericStr(32))
	require.NoError(t, err)
	require.NotNil(t, tokenMaker)

	// initialize uService
	uService := unit.NewService(ur)
	pService := property_service.NewService(pr, ur, lr, ar, rr, authRepo)

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
		ur:         ur,
		pr:         pr,
		tokenMaker: tokenMaker,
		router:     httpServer,
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
