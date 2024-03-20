package http

import (
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/stretchr/testify/require"
	"github.com/user2410/rrms-backend/internal/domain/application"
	"github.com/user2410/rrms-backend/internal/domain/application/asynctask"
	"github.com/user2410/rrms-backend/internal/domain/application/repo"
	chat_repo "github.com/user2410/rrms-backend/internal/domain/chat/repo"
	"github.com/user2410/rrms-backend/internal/domain/listing"
	listing_repo "github.com/user2410/rrms-backend/internal/domain/listing/repo"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
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
	ar repo.Repo, pr property_repo.Repo, ur unit_repo.Repo, lr listing_repo.Repo, cr chat_repo.Repo,
	taskDistributor asynctask.TaskDistributor,
) *server {

	tokenMaker, err := token.NewJWTMaker(random.RandomAlphanumericStr(32))
	require.NoError(t, err)
	require.NotNil(t, tokenMaker)

	// initialize lService
	aService := application.NewService(ar, cr, lr, pr, taskDistributor, nil)
	lService := listing.NewService(lr, pr, nil, "") // NOTE: leave paymentRepo nil for now

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
