package http

import (
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/stretchr/testify/require"
	"github.com/user2410/rrms-backend/internal/domain/listing"
	"github.com/user2410/rrms-backend/internal/domain/listing/repo"
	"github.com/user2410/rrms-backend/internal/domain/property"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	"github.com/user2410/rrms-backend/internal/domain/unit"
	unit_repo "github.com/user2410/rrms-backend/internal/domain/unit/repo"
	"github.com/user2410/rrms-backend/internal/infrastructure/http"
	"github.com/user2410/rrms-backend/internal/utils/random"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

type server struct {
	ur         unit_repo.Repo
	pr         property_repo.Repo
	lr         repo.Repo
	tokenMaker token.Maker
	router     http.Server
}

func newTestServer(t *testing.T, pr property_repo.Repo, ur unit_repo.Repo, lr repo.Repo) *server {

	tokenMaker, err := token.NewJWTMaker(random.RandomAlphanumericStr(32))
	require.NoError(t, err)
	require.NotNil(t, tokenMaker)

	// initialize lService
	uService := unit.NewService(ur)
	pService := property.NewService(pr, ur)
	lService := listing.NewService(lr)

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
	NewAdapter(lService, pService, uService).RegisterServer(httpServer.GetApiRoute(), tokenMaker)

	return &server{
		pr:         pr,
		ur:         ur,
		lr:         lr,
		tokenMaker: tokenMaker,
		router:     httpServer,
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
