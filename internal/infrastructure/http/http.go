package http

import (
	"errors"
	"fmt"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	rcv "github.com/gofiber/fiber/v2/middleware/recover"
)

type Server interface {
	Start(port uint16) error
	GetFibApp() *fiber.App
	GetApiRoute() fiber.Router
	Shutdown() error
}

type server struct {
	fib      *fiber.App
	apiRoute fiber.Router
}

func NewServer(conf fiber.Config) Server {
	return (&server{
		fib: fiber.New(
			conf,
			fiber.Config{
				// Override default error handler
				ErrorHandler: func(ctx *fiber.Ctx, err error) error {
					// Status code defaults to 500
					code := fiber.StatusInternalServerError

					// Retrieve the custom status code if it's a *fiber.Error
					var e *fiber.Error
					if errors.As(err, &e) {
						code = e.Code
					}

					// Send JSON
					err = ctx.Status(code).JSON(fiber.Map{"message": e.Message})
					if err != nil {
						// In case fails
						return ctx.SendStatus(fiber.StatusInternalServerError)
					}

					// Return from handler
					return nil
				},
			},
		),
	}).init()
}

func (s *server) init() Server {
	s.fib.Use(rcv.New())
	s.fib.Use(cors.New())
	s.fib.Use(logger.New())

	s.apiRoute = s.fib.Group("/api")

	return s
}

func (s *server) Start(port uint16) error {
	return s.fib.Listen(fmt.Sprintf(":%d", port))
}

func (s *server) GetFibApp() *fiber.App {
	return s.fib
}

func (s *server) GetApiRoute() fiber.Router {
	return s.apiRoute
}

func (s *server) Shutdown() error {
	return s.fib.Shutdown()
}
