package http

import (
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/stretchr/testify/require"
	"github.com/user2410/rrms-backend/internal/domain/auth"
	"github.com/user2410/rrms-backend/internal/domain/auth/repo"
	"github.com/user2410/rrms-backend/internal/infrastructure/http"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

type server struct {
	r          repo.Repo
	tokenMaker token.Maker
	router     http.Server
}

func newTestServer(t *testing.T, r repo.Repo) *server {
	tokenMaker, err := token.NewJWTMaker("cae1X53au6agHqAOulzCRhgDr0BG52yv")
	require.NoError(t, err)
	require.NotNil(t, tokenMaker)

	// initialize service
	service := auth.NewService(r, tokenMaker, time.Minute, time.Hour)

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
	NewAdapter(service).RegisterServer(httpServer.GetApiRoute(), tokenMaker)

	return &server{
		r:          r,
		tokenMaker: tokenMaker,
		router:     httpServer,
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
