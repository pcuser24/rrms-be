package http

import (
	"fmt"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	cv "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/user2410/rrms-backend/internal/bjwt"
	"github.com/user2410/rrms-backend/internal/infrastructure/repositories/db"
	"log"
	"os"
	"os/signal"
)

type Server interface {
	RegisterApi(dao db.DAO) Server
	GetApiRouter() fiber.Router
	GetAuthRouter() fiber.Router
	Start()
}

type server struct {
	fib    *fiber.App
	bj     bjwt.BJwt
	index  *string
	api    fiber.Router
	auth   fiber.Router
	pub    fiber.Router
	static fiber.Router
}

func NewServer() Server {
	return (&server{
		fib: fiber.New(fiber.Config{}),
		bj:  bjwt.NewBjwt("secret"),
	}).init()
}

func (s *server) Start() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		fmt.Println("Gracefully shutting down...")
		_ = s.fib.Shutdown()
	}()

	if err := s.fib.Listen(fmt.Sprintf(":%d", 8003)); err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Running cleanup tasks...")
}

func (s *server) init() Server {
	s.fib.Use(cv.New())
	s.fib.Use(cors.New())

	s.auth = s.fib.Group("/auth", func(c *fiber.Ctx) error {
		return c.Next()
	})
	s.pub = s.fib.Group("/pub", func(c *fiber.Ctx) error {
		return c.Next()
	})
	s.static = s.fib.Group("/static", func(c *fiber.Ctx) error {
		return c.Next()
	})
	return s
}

func (s *server) GetApiRouter() fiber.Router {
	if s.api == nil {
		s.api = s.fib.Group("/api", func(c *fiber.Ctx) error {
			return c.Next()
		})
		s.api.Use(s.bj.GetJwtHandle())
	}
	return s.api
}

func (s *server) GetAuthRouter() fiber.Router {
	if s.auth == nil {
		s.auth = s.fib.Group("/auth", func(c *fiber.Ctx) error {
			return c.Next()
		})
		s.auth.Use(s.bj.GetJwtHandle())
	}
	return s.auth
}

func (s *server) RegisterApi(dao db.DAO) Server {
	return s
}
