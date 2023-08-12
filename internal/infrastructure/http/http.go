package http

import (
	"fmt"
	"reflect"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	rcv "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
)

type Server interface {
	Start(port uint16) error
	GetFibApp() *fiber.App
	GetApiRoute() *fiber.Router
	Shutdown() error
}

type server struct {
	fib      *fiber.App
	apiRoute *fiber.Router
}

func NewServer(conf fiber.Config) Server {
	return (&server{
		fib: fiber.New(conf),
	}).init()
}

func (s *server) init() Server {
	fiber.SetParserDecoder(fiber.ParserConfig{
		IgnoreUnknownKeys: true,
		ParserType: []fiber.ParserType{
			fiber.ParserType{
				Customtype: uuid.UUID{},
				Converter: func(value string) reflect.Value {
					if v, err := uuid.Parse(value); err == nil {
						return reflect.ValueOf(v)
					}
					return reflect.Value{}
				},
			},
		},
		ZeroEmpty: true,
	})

	s.fib.Use(rcv.New())
	s.fib.Use(cors.New())
	s.fib.Use(logger.New())

	apiRoute := s.fib.Group("/api")
	s.apiRoute = &apiRoute

	return s
}

func (s *server) Start(port uint16) error {
	return s.fib.Listen(fmt.Sprintf(":%d", port))
}

func (s *server) GetFibApp() *fiber.App {
	return s.fib
}

func (s *server) GetApiRoute() *fiber.Router {
	return s.apiRoute
}

func (s *server) Shutdown() error {
	return s.fib.Shutdown()
}
