package http

import (
	"fmt"
	"reflect"
	"runtime/debug"

	"github.com/gofiber/contrib/websocket"
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
	SetWsRoute(route string)
	Shutdown() error
}

type server struct {
	fib      *fiber.App
	apiRoute *fiber.Router
}

func NewServer(conf fiber.Config, corsConf cors.Config) Server {
	return (&server{
		fib: fiber.New(conf),
	}).init(corsConf)
}

func (s *server) init(corsConf cors.Config) Server {
	fiber.SetParserDecoder(fiber.ParserConfig{
		IgnoreUnknownKeys: true,
		ParserType: []fiber.ParserType{
			{
				Customtype: uuid.Nil,
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

	s.fib.Use(cors.New(corsConf))
	s.fib.Use(rcv.New(rcv.Config{
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			debug.PrintStack()
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal server error",
			})
		},
	}))
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

func (s *server) SetWsRoute(route string) {
	s.fib.Use(route, func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
}

func (s *server) Shutdown() error {
	return s.fib.Shutdown()
}
